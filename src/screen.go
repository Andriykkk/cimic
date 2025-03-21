package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/AllenDang/giu"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func onClickMe() {
	fmt.Println("Hello world!")
}

func onImSoCute() {
	cameraSpeed := float32(cameraSpeed) * deltaTime

	cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
	fmt.Println("Im sooooooo cute!!")
}

func loop() {
	giu.SingleWindow().Layout(
		giu.Label("Hello world from giu"),
		giu.Row(
			giu.Button("Click Me").OnClick(onClickMe),
			giu.Button("I'm so cute").OnClick(onImSoCute),
		),
	)
}

func createShaderProgram(vertexSource, fragmentSource string) uint32 {
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatalln("failed to compile vertex shader:", err)
	}
	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatalln("failed to compile fragment shader:", err)
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	gl.UseProgram(shaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderProgram
}

func setUniforms(shaderProgram uint32) {
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, float32(nearClipping), float32(farClipping))
	projectionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
}

func initWindow(vertices []Vertex, faces []Face) {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cimic", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.MakeContextCurrent()
	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyLeftControl || key == glfw.KeyRightControl {
			if action == glfw.Press {
				window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
				lastX, lastY = window.GetCursorPos()
			} else if action == glfw.Release {
				window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
			}
		}

		if key == glfw.KeyEscape && action == glfw.Press {
			window.SetShouldClose(true)
		}
	})
	window.SetCursorPosCallback(mouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL:", err)
	}

	gl.Viewport(0, 0, windowWidth, windowHeight)

	glData := convertToOpenGLData(vertices, faces)

	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(glData)*4, gl.Ptr(glData), gl.STATIC_DRAW)

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// shadow

	vertexShaderSource := `
		#version 410
		layout (location = 0) in vec3 aPos;
		uniform mat4 model;
		uniform mat4 view;
		uniform mat4 projection;
		void main() {
			gl_Position = projection * view * model * vec4(aPos, 1.0);
		} 
		` + "\x00"

	fragmentShaderSource := `
	#version 410
	out vec4 FragColor;
	void main() {
		FragColor = vec4(1.0f, 1.0f, 1.0f, 1.0f);
	}
	` + "\x00"

	wireframeFragmentShaderSource := `
	#version 410
	out vec4 FragColor;
	void main() {
		FragColor = vec4(0.0f, 0.0f, 0.0f, 0.0f);
	}
	` + "\x00"

	// create white triangles
	shaderProgram := createShaderProgram(vertexShaderSource, fragmentShaderSource)
	setUniforms(shaderProgram)

	// create wireframe
	wireframeShaderProgram := createShaderProgram(vertexShaderSource, wireframeFragmentShaderSource)
	setUniforms(wireframeShaderProgram)

	gl.Enable(gl.DEPTH_TEST)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		wnd := giu.NewMasterWindow("Hello world", 400, 200, giu.MasterWindowFlagsNotResizable)
		wnd.Run(loop)
	}()

	for !window.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		processInput(window)

		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		cameraFront = mgl32.Vec3{
			float32(math.Cos(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch))))),
			float32(math.Sin(float64(mgl32.DegToRad(float32(pitch))))),
			float32(math.Sin(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch))))),
		}.Normalize()

		view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)

		// Render solid triangles
		gl.UseProgram(shaderProgram)
		viewUniform := gl.GetUniformLocation(shaderProgram, gl.Str("view\x00"))
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(glData)/3))

		// Render wireframe triangles
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.LineWidth(2.0)
		gl.UseProgram(wireframeShaderProgram)
		wireframeViewUniform := gl.GetUniformLocation(wireframeShaderProgram, gl.Str("view\x00"))
		gl.UniformMatrix4fv(wireframeViewUniform, 1, false, &view[0])
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(glData)/3))
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	wg.Wait()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])
		return 0, fmt.Errorf("shader compilation failed: %s", string(log))
	}

	return shader, nil
}

func mouseCallback(window *glfw.Window, xpos, ypos float64) {
	if window.GetInputMode(glfw.CursorMode) == glfw.CursorDisabled {
		if firstMouse {
			lastX = xpos
			lastY = ypos
			firstMouse = false
		}

		xoffset := xpos - lastX
		yoffset := lastY - ypos
		lastX = xpos
		lastY = ypos

		xoffset *= sensitivity
		yoffset *= sensitivity

		yaw += xoffset
		pitch += yoffset

		if pitch > 89.0 {
			pitch = 89.0
		}
		if pitch < -89.0 {
			pitch = -89.0
		}
	}
}

func processInput(window *glfw.Window) {
	cameraSpeed := float32(cameraSpeed) * deltaTime

	if window.GetKey(glfw.KeyW) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		cameraPos = cameraPos.Add(cameraUp.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press || window.GetKey(glfw.KeyRightShift) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraUp.Mul(cameraSpeed))
	}
}
