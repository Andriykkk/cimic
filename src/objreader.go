package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readObj(filename string) ([]Vertex, []Face, []Vertex, []Vertex, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer file.Close()

	var vertices []Vertex
	var faces []Face
	var normals []Vertex
	var textures []Vertex

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		parts := strings.Fields(line)

		switch parts[0] {
		case "v":
			var x, y, z float64
			fmt.Sscanf(parts[1], "%f", &x)
			fmt.Sscanf(parts[2], "%f", &y)
			fmt.Sscanf(parts[3], "%f", &z)
			vertices = append(vertices, Vertex{X: x, Y: y, Z: z})
		case "vn":
			var x, y, z float64
			fmt.Sscanf(parts[1], "%f", &x)
			fmt.Sscanf(parts[2], "%f", &y)
			fmt.Sscanf(parts[3], "%f", &z)
			normals = append(normals, Vertex{X: x, Y: y, Z: z})
		case "vt":
			var u, v float64
			fmt.Sscanf(parts[1], "%f", &u)
			fmt.Sscanf(parts[2], "%f", &v)
			textures = append(textures, Vertex{X: u, Y: v, Z: 0.0})
		case "f":
			var vertexIndices, textureIndices, normalIndices []int
			for i := 1; i < len(parts); i++ {
				var vertexIdx, textureIdx, normalIdx int
				_, err := fmt.Sscanf(parts[i], "%d/%d/%d", &vertexIdx, &textureIdx, &normalIdx)
				if err != nil {
					_, err = fmt.Sscanf(parts[i], "%d//%d", &vertexIdx, &normalIdx)
					if err != nil {
						_, err = fmt.Sscanf(parts[i], "%d", &vertexIdx)
						if err != nil {
							fmt.Println("Error parsing face:", err)
							continue
						}
					}
					textureIdx = -1
				}

				vertexIndices = append(vertexIndices, vertexIdx-1)
				textureIndices = append(textureIndices, textureIdx-1)
				normalIndices = append(normalIndices, normalIdx-1)
			}

			faces = append(faces, Face{vertexIndices, textureIndices, normalIndices})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}

	return vertices, faces, normals, textures, nil
}

func convertToOpenGLData(vertices []Vertex, faces []Face) []float32 {
	var glData []float32

	for _, face := range faces {
		if len(face.VertexIndices) < 3 {
			continue
		}

		for i := 1; i < len(face.VertexIndices)-1; i++ {
			indices := []int{
				face.VertexIndices[0],
				face.VertexIndices[i],
				face.VertexIndices[i+1],
			}

			for _, idx := range indices {
				v := vertices[idx]
				glData = append(glData, float32(v.X), float32(v.Y), float32(v.Z))
			}
		}
	}

	return glData
}
