package com.water.gdx;

import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.glutils.ShapeRenderer;
import com.badlogic.gdx.math.Vector2;


public class PrimitiveBatch {

    // this constant controls how large the vertices buffer is. Larger buffers will
    // require flushing less often, which can increase performance. However, having
    // buffer that is unnecessarily large will waste memory.
    int DefaultBufferSize = 500;

    // a block of vertices that calling AddVertex will fill. Flush will draw using
    // this array, and will determine how many primitives to draw from
    // positionInBuffer.
    Vector2[] vertices = new Vector2[DefaultBufferSize];

    // keeps track of how many vertices have been added. this value increases until
    // we run out of space in the buffer, at which time Flush is automatically
    // called.
    int positionInBuffer = 0;

    // this value is set by Begin, and is the type of primitives that we are
    // drawing.
    ShapeRenderer.ShapeType primitiveType;

    // how many verts does each of these primitives take up? points are 1,
    // lines are 2, and triangles are 3.
    int numVertsPerPrimitive = 2;

    // hasBegun is flipped to true once Begin is called, and is used to make
    // sure users don't call End before Begin is called.
    boolean hasBegun = false;

    ShapeRenderer sr = new ShapeRenderer();

    // Begin is called to tell the PrimitiveBatch what kind of primitives will be
    // drawn, and to prepare the graphics card to render those primitives.
    public void Begin(ShapeRenderer.ShapeType primitiveType)
    {
        this.primitiveType = primitiveType;
        hasBegun = true;
    }

    // AddVertex is called to add another vertex to be rendered. To draw a point,
    // AddVertex must be called once. for lines, twice, and for triangles 3 times.
    // this function can only be called once begin has been called.
    // if there is not enough room in the vertices buffer, Flush is called
    // automatically.
    public void AddVertex(Vector2 vertex, Color color)
    {
        // are we starting a new primitive? if so, and there will not be enough room
        // for a whole primitive, flush.
        boolean newPrimitive = ((positionInBuffer % numVertsPerPrimitive) == 0);

        if (newPrimitive &&
                (positionInBuffer + numVertsPerPrimitive) >= vertices.length)
        {
            Flush();
        }

        // once we know there's enough room, set the vertex in the buffer,
        // and increase position.
        vertices[positionInBuffer] = vertex;
        positionInBuffer++;
    }

    // End is called once all the primitives have been drawn using AddVertex.
    // it will call Flush to actually submit the draw call to the graphics card, and
    // then tell the basic effect to end.
    public void End()
    {
        // Draw whatever the user wanted us to draw
        Flush();
        hasBegun = false;
    }

    // Flush is called to issue the draw call to the graphics card. Once the draw
    // call is made, positionInBuffer is reset, so that AddVertex can start over
    // at the beginning. End will call this to draw the primitives that the user
    // requested, and AddVertex will call this if there is not enough room in the
    // buffer.
    private void Flush()
    {
        // no work to do
        if (positionInBuffer == 0)
        {
            return;
        }

        sr.setColor(new Color(0, 0.5f, 1.5f, 0.5f));
        sr.begin(this.primitiveType);
        if(this.primitiveType.equals(ShapeRenderer.ShapeType.Filled)) {
            for (int i = 0; i < (vertices.length) - 2; i += 2) {
                sr.triangle(vertices[i].x, vertices[i].y, vertices[i + 1].x, vertices[i+1].y,0,0);
            }
        }
        else if (this.primitiveType.equals(ShapeRenderer.ShapeType.Line)) {
            for (int i = 0; i < (vertices.length) - 2; i += 2) {
                sr.line(vertices[i], vertices[i + 1]);
            }
        }
        sr.end();

        positionInBuffer = 0;
    }
}
