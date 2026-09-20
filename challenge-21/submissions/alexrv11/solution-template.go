package main

import (
	"fmt"
)

func main() {
	// Example sorted array for testing
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

	// Test binary search
	target := 7
	index := BinarySearch(arr, target)
	fmt.Printf("BinarySearch: %d found at index %d\n", target, index)

	// Test recursive binary search
	recursiveIndex := BinarySearchRecursive(arr, target, 0, len(arr)-1)
	fmt.Printf("BinarySearchRecursive: %d found at index %d\n", target, recursiveIndex)

	// Test find insert position
	insertTarget := 8
	insertPos := FindInsertPosition(arr, insertTarget)
	fmt.Printf("FindInsertPosition: %d should be inserted at index %d\n", insertTarget, insertPos)
}

// BinarySearch performs a standard binary search to find the target in the sorted array.
// Returns the index of the target if found, or -1 if not found.
func BinarySearch(arr []int, target int) int {
	// TODO: Implement this function
	start := 0
	last := len(arr) - 1
	for start <= last {
	   pivot := (last - start)/2 + start
	   
	   data := arr[pivot]
	   switch {
	       case data == target:
	       return pivot
	       case data < target:
	        start = pivot + 1
	       default:
	         last = pivot - 1
	   }
	}
	
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, start int, last int) int {
	// TODO: Implement this function
	if len(arr) == 0 || last < start {
	   return -1    
	}
	index := (last - start)/2 + start
	if arr[index] == target {
	    return index
	} else if arr[index] < target {
	    start = index + 1
	} else {
	    last = index - 1
	}
	
	return BinarySearchRecursive(arr, target, start, last)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	index := -1
	last := len(arr) -1
	
	if len(arr) == 0 {
	    return 0
	}
	
	if target <= arr[0] {
	    return 0
	}
	
	if arr[last] < target {
	    return last + 1
	}
	
	if arr[last] == target  {
	    return last
	}
	
	
	
	for i := 0; i < last; i++ {
	    if arr[i] <= target && arr[i+1] >= target {
	        index = i + 1
	        break
	    }
	}
	
	
	return index
}
