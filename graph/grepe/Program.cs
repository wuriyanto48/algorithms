using System;
using System.Collections.Generic;

namespace Grepe
{

    public class User {

        public User(string id, string name)
        {
            Id = id;
            Name = name;
        }
        public string Id {
            get; set;
        }

        public string Name {
            get; set;
        }

        public override string ToString()
        {
            return String.Format($"Id: {Id}, Name: {Name}");
        }
    }
    class Program
    {
        static void Main(string[] args)
        {
            // Graph input representation in matrix
            // 	     | 001 | 002 | 003 | 004 |
            // -------------------------------
            //   001 | 0   | 2   | 0   | 4   |
            //   002 | 0   | 0   | 0   | 3   |
            //   003 | 3   | 2   | 0   | 0   |
            //   004 | 0   | 0   | 1   | 0   |

            var alex = new User("001", "Alex");
            var bony = new User("002", "Bony");
            var deny = new User("003", "Deny");
            var andy = new User("004", "Andy");

            var edges = new List<Graph.Edge<string, User>>();
            edges.Add(new Graph.Edge<string, User>(new Graph.Node<string, User>(alex.Id, alex), new Graph.Node<string, User>(bony.Id, bony), 2));
            edges.Add(new Graph.Edge<string, User>(new Graph.Node<string, User>(alex.Id, alex), new Graph.Node<string, User>(andy.Id, andy), 4));
            edges.Add(new Graph.Edge<string, User>(new Graph.Node<string, User>(bony.Id, bony), new Graph.Node<string, User>(andy.Id, andy), 3));
            edges.Add(new Graph.Edge<string, User>(new Graph.Node<string, User>(deny.Id, deny), new Graph.Node<string, User>(alex.Id, alex), 3));
            edges.Add(new Graph.Edge<string, User>(new Graph.Node<string, User>(deny.Id, deny), new Graph.Node<string, User>(bony.Id, bony), 2));
            edges.Add(new Graph.Edge<string, User>(new Graph.Node<string, User>(andy.Id, andy), new Graph.Node<string, User>(deny.Id, deny), 1));

            var g = new Graph.Graph<string, User>(edges);
            g.ShowEdges();

             Console.WriteLine("-------------------------");
            g.GetNodes(node => {
                Console.WriteLine(node.Value);
            });

            Console.WriteLine("-------------------------");
            // foreach(var b in g.BFS("004")) {
            //     Console.Write($"{b} ---> ");
            // }

            g.DFS("002", node => {
                Console.WriteLine(node);
            });

        //      Console.WriteLine("-------------------------");

        //     Console.WriteLine(g.HasAdjacent("002", "004"));

        //     Console.WriteLine("-------------------------");

        //    g.Neighbors("003", n => {Console.WriteLine(n);});

            Console.WriteLine(g.LengthNodes());
        }
    }
}
