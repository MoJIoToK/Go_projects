package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
)

const bytesInGB = 1024 * 1024 * 1024

// GetDiskUsage retrieves the total, free, and available disk space for a specified path.
func GetDiskUsage(path string) (total, free, available uint64, err error) {
	// Convert the provided path to UTF-16 encoding, required by Windows API.
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to convert path to UTF-16: %w", err)
	}

	// Variables to hold the results from the Windows API call.
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64

	// Call the Windows API function GetDiskFreeSpaceEx to retrieve disk space information.
	err = windows.GetDiskFreeSpaceEx(
		pathPtr,                 // Pointer to the path in UTF-16 format.
		&freeBytesAvailable,     // Pointer to available space for the user.
		&totalNumberOfBytes,     // Pointer to total space on the disk.
		&totalNumberOfFreeBytes, // Pointer to total free space on the disk.
	)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get disk usage: %w", err)
	}

	// Return the retrieved disk space values.
	return totalNumberOfBytes, totalNumberOfFreeBytes, freeBytesAvailable, nil
}

func main() {
	// Specify the disk path for which the usage statistics are needed.
	path := "C:\\"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// Call GetDiskUsage to retrieve disk statistics.
	total, free, available, err := GetDiskUsage(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err) // Print error if the function fails.
		return
	}

	// Convert bytes to gigabytes for easier readability.
	totalGB := float64(total) / float64(bytesInGB)
	freeGB := float64(free) / float64(bytesInGB)
	availableGB := float64(available) / float64(bytesInGB)
	usedGB := totalGB - freeGB // Calculate used space in gigabytes.

	// Print the disk usage statistics in gigabytes.
	fmt.Printf("Disk Usage for Path: %s\n", path)
	fmt.Printf("Total Space: %.2f GB\n", totalGB)
	fmt.Printf("Free Space: %.2f GB\n", freeGB)
	fmt.Printf("Available Space: %.2f GB\n", availableGB)
	fmt.Printf("Used Space: %.2f GB\n", usedGB)
}
