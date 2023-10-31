package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"slashslinging/hasher/models/cluster"
	s "slashslinging/hasher/strategies"
	"slashslinging/hasher/strategies/hashmod"
	"slashslinging/hasher/strategies/hashring"
)

const (
	LOAD = "LOAD"
	INIT = "INIT"
	ADD  = "ADD"
	DEL  = "DEL"
	INFO = "INFO"
	LIST = "LIST"
	PUT  = "PUT"
	GET  = "GET"
	EXIT = "EXIT"
)

const (
	HASHRING = "HRING"
	HASHMOD  = "HMOD"
)

type Command struct {
	ctype string
	cargs []string
}

type Result struct {
	err int
}

type CLI struct {
	clus                 *cluster.Cluster
	distributionStrategy string
	hmStrategy           *hashmod.HashModStrategy
	hrStrategy           *hashring.HashRingStrategy
}

func (cli *CLI) load(path string) {
	file, _ := os.Open(path)
	reader := bufio.NewReader(file)

	line, err := reader.ReadString(byte('\n'))
	for err != io.EOF {
		cmd := cli.read(line)
		result := cli.Eval(cmd)
		if result.err != 0 {
			panic(fmt.Sprintf("error reading/executing command %s", line))
		}
		line, err = reader.ReadString(byte('\n'))
	}
}

func (cli *CLI) eval(cmd Command, ds s.DistributionStrategy) Result {
	switch cmd.ctype {
	case INFO:
		for _, server := range cli.clus.Servers() {
			fmt.Printf("--[Info for server id:%d]--\n", server.Id())
			for _, key := range server.Keys() {
				fmt.Printf("* (key, value): (%s, %s)\n", key, server.Get(key))
			}
		}
		println("==[Strategy Info]==\n")
		ds.Info()
	case LIST:
		idx, _ := strconv.Atoi(cmd.cargs[0])
		fmt.Printf("Printing info for [server %d]\n", idx)
		server := cli.clus.Servers()[idx]
		for _, key := range server.Keys() {
			fmt.Printf("* (key, value): (%s, %s)\n", key, server.Get(key))
		}
	case PUT:
		key := cmd.cargs[0]
		val := cmd.cargs[1]
		idx := ds.GetPartitionIndex(*cli.clus, key)
		cli.clus.Servers()[idx].Put(key, val)
		fmt.Printf("Put (key, value) pair (%s, %s) in server %d\n", key, val, idx)
	case GET:
		key := cmd.cargs[0]
		idx := ds.GetPartitionIndex(*cli.clus, key)
		val := cli.clus.Servers()[idx].Get(key)
		fmt.Printf("Get (%s), value: (%s) from server %d\n", key, val, idx)
	default:
		println("Input not recognized, please try again...")
	}
	return Result{0}
}
func (cli *CLI) Eval(cmd Command) Result {
	switch cmd.ctype {
	case LOAD:
		cli.load(cmd.cargs[0])
	case INIT:
		count, _ := strconv.Atoi(cmd.cargs[0])
		cli.clus = cluster.New(count)
		if cmd.cargs[1] == HASHRING {
			cli.hrStrategy = hashring.New(*cli.clus)
			cli.distributionStrategy = HASHRING
			fmt.Printf("Initializing cluster with %d servers and HASHRING strategy\n", count)
		} else if cmd.cargs[1] == HASHMOD {
			cli.hmStrategy = hashmod.New()
			cli.distributionStrategy = HASHMOD
			fmt.Printf("Initializing cluster with %d servers and HASHMOD strategy\n", count)
		}
	case ADD:
		if cli.distributionStrategy == HASHMOD {
			cli.clus.Add()
			s.Redistribute(cli.hmStrategy, cli.clus, make(map[string]string))
		} else if cli.distributionStrategy == HASHRING {
			servId := cli.clus.Add()
			cli.hrStrategy.Add(servId)
			s.Redistribute(cli.hrStrategy, cli.clus, make(map[string]string))
		}
		println("Added a new server to the cluster")
	case DEL:
		if cli.distributionStrategy == HASHMOD {
			id, _ := strconv.Atoi(cmd.cargs[0])
			server := cli.clus.Del(id)
			s.Redistribute(cli.hmStrategy, cli.clus, server.Map())
		} else if cli.distributionStrategy == HASHRING {
			id, _ := strconv.Atoi(cmd.cargs[0])
			server := cli.clus.Del(id)
			cli.hrStrategy.Del(server.Id())
			s.Redistribute(cli.hrStrategy, cli.clus, server.Map())
			fmt.Printf("Server that was deleted had id: %d\n", server.Id())
		}
	case EXIT:
		return Result{-1}
	default:
		if cli.distributionStrategy == HASHMOD {
			return cli.eval(cmd, *cli.hmStrategy)
		} else if cli.distributionStrategy == HASHRING {
			return cli.eval(cmd, *cli.hrStrategy)
		}
	}
	return Result{0}
}

func (cli *CLI) Print(res Result) {
	// do something in here
}

func (cli *CLI) Repl() int {
	cliHeader()
	defer cliFooter()

	cmd := cli.Read()
	result := cli.Eval(cmd)
	cli.Print(result)
	return result.err
}

func (cli *CLI) Read() Command {
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString(byte('\n'))
	return cli.read(str)
}

func (cli *CLI) read(line string) Command {
	strParts := strings.Split(strings.Trim(line, "\n"), " ")
	return Command{ctype: strings.ToUpper(strParts[0]), cargs: strParts[1:]}
}

func cliHeader() {
	println("Welcome this is a simple (mock key value store cluster")
	println("The purpose of this application is to better understand the advantages to using different")
	println("distribution strategies.")
	println("The commands available to you are as follow: LOAD <str> | INIT <int> <str> | ADD | DEL <int> | INFO | LIST <int> | PUT <str> <str> | GET <str>")
	println("LOAD <url:str> reads commands from and executes them line by line")
	println("INIT <x:int> <s:str>, initializes the cluster with `x` amount of servers and distributes load using `s` strategy")
	println("ADD, adds a server to the cluster")
	println("DEL <id:int>, deletes the server identified by `id` from the cluster")
	println("INFO, Get some information about the cluster")
	println("LIST <id:int>, Get some information about a specific server")
	println("PUT <k:str> <v:str>, stores the k:v pair")
	println("GET <k:str>, if the key exists in the cluster return the value")
	println()
	println()
	println("----")
}

func cliFooter() {
	println("----")
	println()
	println()
}

func main() {
	cli := new(CLI)
	for {
		exit := cli.Repl()

		if exit != 0 {
			println("exiting...")
			break
		}
	}
}
