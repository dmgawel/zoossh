// Tests functions from "generic.go".

package zoossh

import (
	"os"
	"testing"
)

// Test the function ParseUnknownFile().
func TestParseUnknownFile(t *testing.T) {

	_, err := ParseUnknownFile("/dev/zero")
	if err == nil {
		t.Errorf("ParseUnknownFile() failed to reject /dev/zero.")
	}

	// Only run this test if the consensus file is there.
	if _, err = os.Stat(consensusFile); err == nil {
		_, err := ParseUnknownFile(consensusFile)
		if err != nil {
			t.Errorf("ParseUnknownFile() failed to parse %s.", consensusFile)
		}
	}
}

func TestInterfaces(t *testing.T) {

	testFingerprint := Fingerprint("9695DFC35FFEB861329B9F1AB04C46397020CE31")

	if _, err := os.Stat(consensusFile); err != nil {
		return
	}

	if _, err := os.Stat(serverDescriptorFile); err != nil {
		return
	}

	consensus, err := ParseUnknownFile(consensusFile)
	if err != nil {
		t.Fatal(err)
	}

	descriptors, err := ParseUnknownFile(serverDescriptorFile)
	if err != nil {
		t.Fatal(err)
	}

	// Test the GetObject() function.
	obj1, found := consensus.GetObject(testFingerprint)
	if found != true {
		t.Error("Could not find existing router status in consensus.")
	}

	obj2, found := descriptors.GetObject(testFingerprint)
	if found != true {
		t.Error("Could not find existing descriptor in descriptor set.")
	}

	// Test the GetFingerprint() function.
	if obj1.GetFingerprint() != testFingerprint {
		t.Error("Failed to retrieve correct fingerprint.")
	}

	if obj2.GetFingerprint() != testFingerprint {
		t.Error("Failed to retrieve correct fingerprint.")
	}

	// Test the Length() function.
	if consensus.Length() != 6840 {
		t.Error("Failed to determine consensus length.")
	}

	if descriptors.Length() != 763 {
		t.Error("Failed to determine descriptor set length.")
	}

	// Test the Iterate() function.
	counter := 0
	for _ = range consensus.Iterate() {
		counter += 1
	}
	if counter != consensus.Length() {
		t.Error("Failed to iterate over all router statuses in consensus.")
	}

	counter = 0
	for _ = range descriptors.Iterate() {
		counter += 1
	}
	if counter != descriptors.Length() {
		t.Error("Failed to iterate over all descriptors in descriptor set.")
	}

	// Test the Merge() function.
	prevLength := consensus.Length()
	consensus.Merge(consensus)
	if consensus.Length() != prevLength {
		t.Error("Consensus merge with itself caused unexpected length.")
	}

	prevLength = descriptors.Length()
	descriptors.Merge(descriptors)
	if descriptors.Length() != prevLength {
		t.Error("Descriptors merge with itself caused unexpected length.")
	}
}
