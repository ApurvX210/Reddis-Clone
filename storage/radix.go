package storage

// Will implement the radix tree data structure for optimizing ttl workflow

type Radix struct{
	prefix	string
	child	map[string] *Radix
	end		bool
}

