using System;
using System.Collections.Generic;
using System.Linq;

namespace Grepe.Graph {

    public struct Node<K, T> : IComparable {

        private K key;
        private T value;

        public Node(K key, T value) {
            this.key = key;
            this.value = value;
        }

        public K Key {
            get => key;
        }

        public T Value {
            get => value;
        }

        public int CompareTo(object obj)
        {
            throw new NotImplementedException();
        }

        public override string ToString()
        {
            return String.Format($"Node {key}");
        }

    }

    public struct Edge<K, T> {
        private Node<K, T> source;
        private Node<K, T> destination;
        private float weight;

        public Edge(Node<K, T> source, Node<K, T> destination, int weight) {
            this.source = source;
            this.destination = destination;
            this.weight = weight;
        }

        public Node<K, T> Source {
            get => source;
        }

        public Node<K, T> Destination {
            get => destination;
        }

        public float Weight {
            get => weight;
        }

        public override string ToString()
        {
            return String.Format($"{source.Key} -> {destination.Key}");
        }
    }

    public class Graph<K, T> {

        private Dictionary<K, List<Edge<K, T>>> adjacencyList;
        private HashSet<Node<K, T>> nodes;

        public Graph(List<Edge<K, T>> edges) {
            adjacencyList = new Dictionary<K, List<Edge<K, T>>>();
            nodes = new HashSet<Node<K, T>>();

            foreach (var e in edges)
            {
                adjacencyList[e.Source.Key] = new List<Edge<K, T>>();
            }

            foreach (var e in edges)
            {
                adjacencyList[e.Source.Key].Add(e);
                nodes.Add(e.Source);
            }

        }

        public void AddEdge(Edge<K, T> e)
        {
            if (!this.adjacencyList.ContainsKey(e.Source.Key)) {
                this.adjacencyList[e.Source.Key] = new List<Edge<K, T>>();
                this.nodes.Add(e.Source);
            }

            this.adjacencyList[e.Source.Key].Add(e);
        }

        public void AddNode(Node<K, T> n)
        {
            if (!this.adjacencyList.ContainsKey(n.Key))
            {
                this.adjacencyList[n.Key] = new List<Edge<K, T>>();
                this.nodes.Add(n);
            }
        }

        public void GetNodes(Action<Node<K, T>> action)
        {
            foreach(var node in nodes)
            {
                action(node);
            }
        }

        public int LengthNodes()
        {
            return nodes.Count;
        }

        public void GetEdges(Action<List<Edge<K, T>>> action)
        {
            foreach(KeyValuePair<K, List<Edge<K, T>>> entry in adjacencyList)
            {
                action(entry.Value);
            }
        }

        public void ShowEdges()
        {
            foreach(KeyValuePair<K, List<Edge<K, T>>> entry in adjacencyList)
            {
                Console.WriteLine($"Adjacency list of vertex {entry.Key}");
                Console.Write("head ");
                foreach(var edge in entry.Value)
                {
                    Console.Write($"--> {edge.Destination.Key} .(weight {edge.Weight}). ");
                }

                Console.Write("\n");
            }
        }

        public bool HasAdjacent(K from, K to)
        {
            if (!this.adjacencyList.ContainsKey(from))
            {
                return false;
            }

            var adjacency = this.adjacencyList[from];
            return adjacency.Exists(o => o.Destination.Key.Equals(to));
        }

        public void Neighbors(K from, Action<Node<K, T>> action) {
            if (!this.adjacencyList.ContainsKey(from))
            {
                return;
            }

            var adjacency = this.adjacencyList[from];
            adjacency.ForEach(e => action(e.Destination));
        }

        // Dfs an implementation of Depth-first search
        public void DFS(K start, Action<K> action)
        {
            var visited = new Dictionary<K, bool>();
            foreach(KeyValuePair<K, List<Edge<K, T>>> entry in adjacencyList)
            {
                visited[entry.Key] = false;
            }

            DFS<K, T>(adjacencyList, start, visited, action);
        }
        
        private static void DFS<I, V>(Dictionary<I, List<Edge<I, V>>> data, I start, Dictionary<I, bool> visited, Action<I> action)
        {
            visited[start] = true;
            //Console.Write($"{start} ---> ");
            action(start);
            foreach(KeyValuePair<I, List<Edge<I, V>>> entry in data)
            {
                if (!visited[entry.Key])
                {
                    visited[entry.Key] = true;
                    DFS(data, entry.Key, visited, action);
                }
            }
        }

        // BFS an implementation of Breadth-first Traversal
        public IEnumerable<K> BFS(K start)
        {
            var queue = new LinkedList<K>();
            var visited = new Dictionary<K, bool>();
            foreach(KeyValuePair<K, List<Edge<K, T>>> entry in adjacencyList)
            {
                visited[entry.Key] = false;
            }

            visited[start] = true;
            queue.AddLast(start);

            while (queue.Any())
            {
                var key = queue.First();
                //Console.Write($"{key} ---> ");
                yield return key;

                // remove value
                queue.RemoveFirst();

                foreach(KeyValuePair<K, List<Edge<K, T>>> entry in adjacencyList)
                {
                    if (!visited[entry.Key])
                    {
                        visited[entry.Key] = true;

                        // push value to queue
                        queue.AddLast(entry.Key);
                    }
                }
               
            }
        }

    }
}