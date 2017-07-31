// Package retag provides an ability to change tags of structures' fields in runtime
// without copying of the data. It may be helpful in next cases:
//
//  - Automatic tags generation;
//  - Different views of the one data;
//  - Fixing of leaky abstractions with minimal boilerplate code
//    when application has layers of abstractions and model is
//    separated from storages and presentation layers.
//
// Features:
//  - No memory allocations (for cached types);
//  - Fast converting (lookup in table and pointer creation for cached types).
//
// The package requires go1.7+.
// Warning! The package probably will not work in go1.9+ with anonymous and unexported fields.
// The package is still experimental and subject to change.
//
package retag
