package core

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"testing"
)

var testInput = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam ullamcorper tempus sagittis. Vestibulum rutrum nisi finibus mollis rhoncus. Aliquam erat volutpat. Sed porttitor ex eget elementum lobortis. Nullam ut dolor at sapien vestibulum fermentum ac eget ante. Mauris convallis dui eu laoreet bibendum. Mauris ante arcu, porta et lacus id, cursus sodales justo. Mauris sed nisi vehicula, lacinia ligula mattis, semper arcu. Pellentesque in molestie diam, non pretium lacus. Integer ante augue, porttitor a lobortis in, pharetra nec diam. Sed scelerisque purus a neque faucibus, sit amet commodo nisi tincidunt. Etiam sapien nibh, venenatis quis convallis sed, scelerisque sit amet ante. Integer urna odio, suscipit sit amet sapien nec, convallis consectetur massa. Aliquam malesuada lectus justo, vitae sollicitudin mauris mattis sit amet. Mauris condimentum iaculis interdum. Aliquam iaculis leo mauris, sit amet dictum ex ultricies in. Aliquam a enim vel neque porta eleifend ut ac quam. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec commodo massa et varius ultrices. Duis varius lobortis ex, ac eleifend odio venenatis ut. Phasellus non luctus eros. Sed sodales lectus id odio porta porta. Curabitur id lacinia quam. Etiam finibus nisi quis velit dapibus auctor. Integer sollicitudin, felis et ornare iaculis, diam risus hendrerit lacus, vel lacinia sem dui vitae erat.")

func TestCompressDecompressBytes(t *testing.T) {
	compressed := CompressBytes(testInput)
	t.Log("compressed: " + string(compressed) + "\n\n")

	decompressed := DecompressBytes(compressed)
	t.Log("decompressed: " + string(decompressed) + "\n\n")
}
