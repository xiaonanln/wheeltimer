# go-wheeltimer
wheel timer that works just like time.After, time.Tick, time.NewTicker（implementing）

Steps: 
1. Read paper 《Hashed and Hierarchical Timing Wheels: Efficient
Data Structures for Implementing a Timer Facility》
2. Implement it

The wheeltimer implementation works now, but it is not as efficient as the standard time library. 
So it is not meaningful for production for now.