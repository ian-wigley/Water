package com.water.gdx;

import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.OrthographicCamera;
import com.badlogic.gdx.graphics.Texture;
import com.badlogic.gdx.graphics.g2d.SpriteBatch;
import com.badlogic.gdx.graphics.glutils.ShapeRenderer;
import com.badlogic.gdx.math.Vector2;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;

public class Water {

    class WaterColumn
    {
        public float TargetHeight;
        public float Height;
        public float Speed;

        public WaterColumn(int Height, int TargetHeight, int Speed){
            this.Height = Height;
            this.TargetHeight = TargetHeight;
            this.Speed = Speed;
        }

        public void Update(float dampening, float tension)
        {
            float x = TargetHeight - Height;
            Speed += tension * x - Speed * dampening;
            Height += Speed;
        }
    }

    class Particle
    {
        public Vector2 Position;
        public Vector2 Velocity;
        public float Orientation;

        public Particle(Vector2 position, Vector2 velocity, float orientation)
        {
            Position = position;
            Velocity = velocity;
            Orientation = orientation;
        }
    }

    PrimitiveBatch pb;
    WaterColumn[] columns = new WaterColumn[201];
    static Random rand = new Random();

    public float Tension = 0.025f;
    public float Dampening = 0.025f;
    public float Spread = 0.25f;

    SpriteBatch spriteBatch;
    Texture particleTexture;
    List<Particle> particles = new ArrayList<Particle>();

    private OrthographicCamera cam;
    public Water(Texture particleTexture)
    {
        cam = new OrthographicCamera(Gdx.graphics.getWidth(), Gdx.graphics.getHeight());
        cam.setToOrtho(true, Gdx.graphics.getWidth(), Gdx.graphics.getHeight());

        pb = new PrimitiveBatch();
        this.particleTexture = particleTexture;
        spriteBatch = new SpriteBatch();
        for (int i = 0; i < columns.length; i++)
        {
            columns[i] = new WaterColumn(240,240,0);
        }
    }

    // Returns the height of the water at a given x coordinate.
    public float GetHeight(float x)
    {
        if (x < 0 || x > 800) {
            return 240;
        }

        return columns[(int)(x / getScale())].Height;
    }

    void UpdateParticle(Particle particle)
    {
        float Gravity = 0.3f;
        particle.Velocity.y += Gravity;
        particle.Position.x += particle.Velocity.x;
        particle.Position.y += particle.Velocity.y;
        particle.Orientation = GetAngle(particle.Velocity);
    }

    private float clamp(float value, int MIN_VALUE, int MAX_VALUE){
        return value > MAX_VALUE ? MAX_VALUE : value < MIN_VALUE ? MIN_VALUE : value;
    }

    public void Splash(float xPosition, float speed)
    {
        int index = (int)clamp(xPosition / getScale(), 0, columns.length - 1);



        index = 100;





        for (int i = Math.max(0, index); i < Math.min(columns.length - 1, index + 1); i++)
            columns[index].Speed = speed;

//        CreateSplashParticles(xPosition, speed);
    }

    private float getScale() {
        return Gdx.graphics.getWidth() / (columns.length - 1f);
    }

    private void CreateSplashParticles(float xPosition, float speed)
    {
        float y = GetHeight(xPosition);

        if (speed > 120)
        {
            for (int i = 0; i < speed / 8; i++)
            {
                Vector2 temp = GetRandomVector2(40);
                Vector2 pos = new Vector2(xPosition+temp.x, y+temp.y);
                Vector2 vel = FromPolar((float) Math.toRadians(GetRandomFloat(-150.0f, -30.0f)), GetRandomFloat(0, 0.5f * (float)Math.sqrt(speed)));
                CreateParticle(pos, vel);
            }
        }
    }

    private void CreateParticle(Vector2 pos, Vector2 velocity)
    {
        particles.add(new Particle(pos, velocity, 0));
    }

    private Vector2 FromPolar(float angle, float magnitude)
    {
        return new Vector2(((float)Math.cos(angle)) * magnitude, ((float)Math.sin(angle)) * magnitude);
    }

    private float GetRandomFloat(float min, float max)
    {
        return (float)rand.nextDouble() * (max - min) + min;
    }

    private Vector2 GetRandomVector2(float maxLength)
    {
        return FromPolar(GetRandomFloat((float)-Math.PI, (float)Math.PI), GetRandomFloat(0, maxLength));
    }

    private float GetAngle(Vector2 vector)
    {
        return (float)Math.atan2(vector.y, vector.x);
    }

    public void Update()
    {
        for (WaterColumn column : columns) column.Update(Dampening, Tension);

        float[] lDeltas = new float[columns.length];
        float[] rDeltas = new float[columns.length];

        // do some passes where columns pull on their neighbours
        for (int j = 0; j < 8; j++)
        {
            for (int i = 0; i < columns.length; i++)
            {
                if (i > 0)
                {
                    lDeltas[i] = Spread * (columns[i].Height - columns[i - 1].Height);
                    columns[i - 1].Speed += lDeltas[i];
                }
                if (i < columns.length - 1)
                {
                    rDeltas[i] = Spread * (columns[i].Height - columns[i + 1].Height);
                    columns[i + 1].Speed += rDeltas[i];
                }
            }

            for (int i = 0; i < columns.length; i++)
            {
                if (i > 0)
                    columns[i - 1].Height += lDeltas[i];
                if (i < columns.length - 1)
                    columns[i + 1].Height += rDeltas[i];
            }
        }

        for (Particle particle : particles) {
            UpdateParticle(particle);
        }
    }

    public void DrawToRenderTargets()
    {
        cam.update();
        spriteBatch.setProjectionMatrix(cam.combined);
        spriteBatch.begin();
        for (Particle particle : particles)
        {
            spriteBatch.draw(particleTexture, particle.Position.x, particle.Position.y, particleTexture.getWidth(), particleTexture.getHeight(), 0, 0,
                    particleTexture.getWidth(),	particleTexture.getHeight(), false, true);
        }
        spriteBatch.end();

//        // draw a gradient above the water so the metaballs will fuse with the water's surface.
//        pb.Begin(ShapeRenderer.ShapeType.Filled);
//
//		float thickness = 20;
//        float scale = getScale();
//        for (int i = 1; i < columns.length; i++)
//        {
//            Vector2 p1 = new Vector2((i - 1) * scale, columns[i - 1].Height);
//            Vector2 p2 = new Vector2(i * scale, columns[i].Height);
//            Vector2 p3 = new Vector2(p1.x, p1.y - thickness);
//            Vector2 p4 = new Vector2(p2.x, p2.y - thickness);
//            pb.AddVertex(p2, Color.WHITE);
//            pb.AddVertex(p1, Color.WHITE);
//            pb.AddVertex(p3, Color.WHITE);
//            pb.AddVertex(p3, Color.WHITE);
//            pb.AddVertex(p4, Color.WHITE);
//            pb.AddVertex(p2, Color.WHITE);
//        }
//        pb.End();
    }

    public void Draw()
    {
        Color lightBlue = new Color(0.2f, 0.5f, 1f, 0f);

        // draw the waves
        pb.Begin(ShapeRenderer.ShapeType.Filled);
        Color midnightBlue = new Color(0, 15, 40, 0);// * 0.9f;

        float bottom = -Gdx.graphics.getHeight();
        float scale = getScale();
        for (int i = 1; i < columns.length; i++)
        {
            Vector2 p1 = new Vector2((i - 1) * scale, columns[i - 1].Height);
            Vector2 p2 = new Vector2(i * scale, columns[i].Height);
            Vector2 p3 = new Vector2(p2.x, bottom);
            Vector2 p4 = new Vector2(p1.x, bottom);

            pb.AddVertex(p1, lightBlue);
            pb.AddVertex(p2, lightBlue);
            pb.AddVertex(p3, midnightBlue);

            pb.AddVertex(p1, lightBlue);
            pb.AddVertex(p3, midnightBlue);
            pb.AddVertex(p4, midnightBlue);
        }

        pb.End();
    }
}
