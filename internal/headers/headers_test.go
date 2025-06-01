package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header without sep space
	headers = NewHeaders()
	data = []byte("Host:localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 22, n)
	assert.False(t, done)

	// Test: Valid header with extra spaces
	headers = NewHeaders()
	data = []byte("         Host:     localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 39, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("         Host : localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nContent-Type: application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, 23, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, 2, len(headers))
	assert.Equal(t, 32, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, 23, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: field name with capital letters
	headers = NewHeaders()
	data = []byte("HOsTABcDefGhIJkL: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["hostabcdefghijkl"])
	assert.Equal(t, 35, n)
	assert.False(t, done)

	// Test: field name with valid special characters
	headers = NewHeaders()
	data = []byte("Ho$t: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["ho$t"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: field name with invalid character
	headers = NewHeaders()
	data = []byte("Ho@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: multiple values for same header
	headers = NewHeaders()
	data = []byte("Set-Person: p1\r\nContent-Type: application/json\r\nSet-Person: p2\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "p1", headers["set-person"])
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, 16, n)
	assert.False(t, done)
	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "p1", headers["set-person"])
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, 2, len(headers))
	assert.Equal(t, 32, n)
	assert.False(t, done)
	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "p1, p2", headers["set-person"])
	assert.Equal(t, 2, len(headers))
	assert.Equal(t, 16, n)
	assert.False(t, done)
}
