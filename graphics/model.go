package graphics

import "C"

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/PetrusJPrinsloo/learnopengl/graphics/assimp"
	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	texturesLoaded  map[string]Texture
	wg              sync.WaitGroup
	Meshes          []Mesh
	GammaCorrection bool
	BasePath        string
	FileName        string
	GobName         string
}

func NewModel(b, f string, g bool) (Model, error) {
	t := strings.Split(f, ".")
	gf := t[0] + ".gob"
	m := Model{
		BasePath:        b,
		FileName:        f,
		GobName:         gf,
		GammaCorrection: g,
	}
	m.texturesLoaded = make(map[string]Texture)
	gobFile := b + gf
	if _, err := os.Stat(gobFile); os.IsNotExist(err) {
		err := m.loadModel()
		m.Export()
		return m, err
	}
	err := m.Import()
	return m, err
}

func (m *Model) Draw(shader uint32) {
	for i := 0; i < len(m.Meshes); i++ {
		m.Meshes[i].draw(shader)
	}
}

func (m *Model) Export() error {
	// export a gob file
	f := m.BasePath + m.GobName

	dataFile, err := os.Create(f)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(m)

	return nil
}

func (m *Model) Import() error {
	f := m.BasePath + m.GobName
	dataFile, err := os.Open(f)

	if err != nil {
		return err
	}
	defer dataFile.Close()

	dataDecoder := gob.NewDecoder(dataFile)
	if err := dataDecoder.Decode(&m); err != nil {
		return err
	}

	fmt.Printf("Creating model from gob file: %s\n", f)
	m.initGL()
	return nil
}

func (m *Model) Dispose() {
	for i := 0; i < len(m.Meshes); i++ {
		gl.DeleteVertexArrays(1, &m.Meshes[i].vao)
		gl.DeleteBuffers(1, &m.Meshes[i].vbo)
		gl.DeleteBuffers(1, &m.Meshes[i].ebo)
	}
}

// Loads a model with supported ASSIMP extensions from file and stores the resulting meshes in the meshes vector.
func (m *Model) loadModel() error {
	// Read file via ASSIMP
	path := m.BasePath + m.FileName
	scene := assimp.ImportFile(path, uint(
		assimp.Process_Triangulate|assimp.Process_FlipUVs))

	// Check for errors
	if scene.Flags()&assimp.SceneFlags_Incomplete != 0 { // if is Not Zero
		fmt.Println("ERROR::ASSIMP:: %s\n", scene.Flags())
		return errors.New("shit failed")
	}

	// Process ASSIMP's root node recursively
	m.processNode(scene.RootNode(), scene)
	m.wg.Wait()
	m.initGL()
	return nil
}

func (m *Model) initGL() {
	// using a for loop with a range doesnt work here?!
	// also making a temp var inside the loop doesnt work either?!
	for i := 0; i < len(m.Meshes); i++ {
		for j := 0; j < len(m.Meshes[i].Textures); j++ {
			if val, ok := m.texturesLoaded[m.Meshes[i].Textures[j].Path]; ok {
				m.Meshes[i].Textures[j].id = val.id
			} else {
				m.Meshes[i].Textures[j].id = m.textureFromFile(m.Meshes[i].Textures[j].Path)
				m.texturesLoaded[m.Meshes[i].Textures[j].Path] = m.Meshes[i].Textures[j]
			}
		}
		m.Meshes[i].setup()
	}
}

func (m *Model) processNode(n *assimp.Node, s *assimp.Scene) {
	// Process each mesh located at the current node
	m.wg.Add(n.NumMeshes() + n.NumChildren())

	for i := 0; i < n.NumMeshes(); i++ {
		// The node object only contains indices to index the actual objects in the scene.
		// The scene contains all the data, node is just to keep stuff organized (like relations between nodes).
		go func(index int) {
			defer m.wg.Done()
			mesh := s.Meshes()[n.Meshes()[index]]
			ms := m.processMesh(mesh, s)
			m.Meshes = append(m.Meshes, ms)
		}(i)

	}

	// After we've processed all of the meshes (if any) we then recursively process each of the children nodes
	c := n.Children()
	for j := 0; j < len(c); j++ {
		go func(n *assimp.Node, s *assimp.Scene) {
			defer m.wg.Done()
			m.processNode(n, s)
		}(c[j], s)
	}
}

func (m *Model) processMeshVertices(mesh *assimp.Mesh) []Vertex {
	// Walk through each of the mesh's vertices
	vertices := []Vertex{}

	positions := mesh.Vertices()

	normals := mesh.Normals()
	useNormals := len(normals) > 0

	tex := mesh.TextureCoords(0)
	useTex := true
	if tex == nil {
		useTex = false
	}

	tangents := mesh.Tangents()
	useTangents := len(tangents) > 0

	bitangents := mesh.Bitangents()
	useBitTangents := len(bitangents) > 0

	for i := 0; i < mesh.NumVertices(); i++ {
		// We declare a placeholder vector since assimp uses its own vector class that
		// doesn't directly convert to glm's vec3 class so we transfer the data to this placeholder glm::vec3 first.
		vertex := Vertex{}

		// Positions
		vertex.Position = mgl.Vec3{positions[i].X(), positions[i].Y(), positions[i].Z()}

		// Normals
		if useNormals {
			vertex.Normal = mgl.Vec3{normals[i].X(), normals[i].Y(), normals[i].Z()}
			//n.WriteString(fmt.Sprintf("[%f, %f, %f]\n", tmp[i].X(), tmp[i].Y(), tmp[i].Z()))
		}

		// Texture Coordinates
		if useTex {
			// Does the mesh contain texture coordinates?
			// A vertex can contain up to 8 different texture coordinates. We thus make the assumption that we won't
			// use models where a vertex can have multiple texture coordinates so we always take the first set (0).
			vertex.TexCoords = mgl.Vec2{tex[i].X(), tex[i].Y()}
		} else {
			vertex.TexCoords = mgl.Vec2{0.0, 0.0}
		}

		// Tangent
		if useTangents {
			vertex.Tangent = mgl.Vec3{tangents[i].X(), tangents[i].Y(), tangents[i].Z()}
		}

		// Bitangent
		if useBitTangents {
			vertex.Bitangent = mgl.Vec3{bitangents[i].X(), bitangents[i].Y(), bitangents[i].Z()}
		}

		vertices = append(vertices, vertex)
	}

	return vertices
}

func (m *Model) processMeshIndices(mesh *assimp.Mesh) []uint32 {
	indices := []uint32{}
	// Now wak through each of the mesh's faces (a face is a mesh its triangle) and retrieve the corresponding vertex indices.
	for i := 0; i < mesh.NumFaces(); i++ {
		face := mesh.Faces()[i]
		// Retrieve all indices of the face and store them in the indices vector
		indices = append(indices, face.CopyIndices()...)
	}
	return indices
}

func (m *Model) processMeshTextures(mesh *assimp.Mesh, s *assimp.Scene) []Texture {
	textures := []Texture{}
	// Process materials
	if mesh.MaterialIndex() >= 0 {
		material := s.Materials()[mesh.MaterialIndex()]

		// We assume a convention for sampler names in the shaders. Each diffuse texture should be named
		// as 'texture_diffuseN' where N is a sequential number ranging from 1 to MAX_SAMPLER_NUMBER.
		// Same applies to other texture as the following list summarizes:
		// Diffuse: texture_diffuseN
		// Specular: texture_specularN
		// Normal: texture_normalN

		// 1. Diffuse maps
		diffuseMaps := m.loadMaterialTextures(material, assimp.TextureMapping_Diffuse, "texture_diffuse")
		textures = append(textures, diffuseMaps...)
		// 2. Specular maps
		specularMaps := m.loadMaterialTextures(material, assimp.TextureMapping_Specular, "texture_specular")
		textures = append(textures, specularMaps...)
		// 3. Normal maps
		normalMaps := m.loadMaterialTextures(material, assimp.TextureMapping_Height, "texture_normal")
		textures = append(textures, normalMaps...)
		// 4. Height maps
		heightMaps := m.loadMaterialTextures(material, assimp.TextureMapping_Ambient, "texture_height")
		textures = append(textures, heightMaps...)
	}
	return textures
}

func (m *Model) processMesh(ms *assimp.Mesh, s *assimp.Scene) Mesh {
	// Return a mesh object created from the extracted mesh data
	return NewMesh(
		m.processMeshVertices(ms),
		m.processMeshIndices(ms),
		m.processMeshTextures(ms, s))
}

func (m *Model) loadMaterialTextures(ms *assimp.Material, tm assimp.TextureMapping, tt string) []Texture {
	textureType := assimp.TextureType(tm)
	textureCount := ms.GetMaterialTextureCount(textureType)
	result := []Texture{}

	for i := 0; i < textureCount; i++ {
		file, _, _, _, _, _, _, _ := ms.GetMaterialTexture(textureType, 0)
		filename := m.BasePath + file
		texture := Texture{id: 0, TextureType: tt, Path: filename}
		result = append(result, texture)

		//if val, ok := m.texturesLoaded[filename]; ok {
		//	result = append(result, val)
		//} else {
		//	texId := m.textureFromFile(filename)
		//	texture := Texture{id: texId, TextureType: tt, Path: filename}
		//	result = append(result, texture)
		//	m.texturesLoaded[filename] = texture
		//}
	}
	return result
}

func (m *Model) textureFromFile(f string) uint32 {
	//Generate texture ID and load texture data
	if tex, err := NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR_MIPMAP_LINEAR, gl.LINEAR, f); err != nil {
		panic(err)
	} else {
		return tex
	}
}
