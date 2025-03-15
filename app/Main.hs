module Main (main) where

import Graphics.UI.GLUT
import Data.IORef

moveX, moveY, moveZ :: GLfloat
moveX = 0
moveY = 0
moveZ = 70

initMovement :: IO (IORef GLfloat, IORef GLfloat, IORef GLfloat)
initMovement = do
    moveX <- newIORef 0
    moveY <- newIORef 0
    moveZ <- newIORef 70
    return (moveX, moveY, moveZ)

initOpenGL :: IO ()
initOpenGL = do
    clearColor $= Color4 0 0 0 1
    matrixMode $= Projection
    loadIdentity
    ortho (-100) 100 (-100) 100 (-100) 100
    matrixMode $= Modelview 0

display :: IORef GLfloat -> IORef GLfloat -> IORef GLfloat -> IO ()
display moveXRef moveYRef moveZRef = do
    clear [ColorBuffer]
    loadIdentity

    x <- readIORef moveXRef
    y <- readIORef moveYRef
    z <- readIORef moveZRef

    translate (Vector3 x y z)

    color (Color3 1.0 0.0 (0.0 :: GLfloat))

    renderPrimitive Triangles $ do
        vertex (Vertex3 (10 :: GLfloat) (10 :: GLfloat) (0 :: GLfloat))
        vertex (Vertex3 (90 :: GLfloat) (10 :: GLfloat) (0 :: GLfloat))
        vertex (Vertex3 (50 :: GLfloat) (90 :: GLfloat) (0 :: GLfloat))
    
    flush

reshape :: ReshapeCallback
reshape size = do
    viewport $= (Position 0 0, size)
    postRedisplay Nothing

idle :: IORef GLfloat -> IORef GLfloat -> IORef GLfloat -> IdleCallback
idle moveXRef moveYRef moveZRef = do
    postRedisplay Nothing

keyboard :: IORef GLfloat -> IORef GLfloat -> IORef GLfloat -> KeyboardCallback
keyboard moveXRef moveYRef moveZRef key _ = case key of
    'w' -> modifyIORef moveZRef (+ 1)  
    's' -> modifyIORef moveZRef (subtract 1)  
    'a' -> modifyIORef moveXRef (subtract 1)  
    'd' -> modifyIORef moveXRef (+ 1) 
    'r' -> modifyIORef moveYRef (+ 1) 
    'f' -> modifyIORef moveYRef (subtract 1)  
    _   -> return ()  

main :: IO ()
main = do
  (_progName, _args) <- getArgsAndInitialize
  _ <- createWindow "GLUT Triangle"
  initOpenGL

  (moveXRef, moveYRef, moveZRef) <- initMovement

  displayCallback $= display moveXRef moveYRef moveZRef
  reshapeCallback $= Just reshape
  idleCallback $= Just (idle moveXRef moveYRef moveZRef)
  keyboardCallback $= Just (keyboard moveXRef moveYRef moveZRef)

  mainLoop
