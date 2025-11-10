package com.water.gdx;

import com.badlogic.gdx.ApplicationAdapter;
import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.Input;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.GL20;
import com.badlogic.gdx.graphics.OrthographicCamera;
import com.badlogic.gdx.graphics.Texture;
import com.badlogic.gdx.graphics.g2d.BitmapFont;
import com.badlogic.gdx.graphics.g2d.SpriteBatch;
import com.badlogic.gdx.math.Vector2;

public class Game extends ApplicationAdapter {
    SpriteBatch batch;
    Texture particleImage;
    Texture backgroundImage;
    Texture rockImage;
    Texture background;
    Water water;
    Rock rock;
    private BitmapFont font;
    private OrthographicCamera cam;
    int x;

    @Override
    public void create() {
        x = 0;
        batch = new SpriteBatch();
        particleImage = new Texture("metaparticle.png");
        backgroundImage = new Texture("sky.png");
        rockImage = new Texture("rock.png");
        background = new Texture("Grounds.png");
        water = new Water(particleImage);

        font = new BitmapFont(true);
        font.setColor(Color.WHITE);

        cam = new OrthographicCamera(Gdx.graphics.getWidth(), Gdx.graphics.getHeight());
        cam.setToOrtho(true, Gdx.graphics.getWidth(), Gdx.graphics.getHeight());
    }

    @Override
    public void render() {
        if (rock != null) {
            if (rock.Position.y < 240 && rock.Position.y + rock.Velocity.y >= 240) {
                water.Splash(rock.Position.x, rock.Velocity.y * rock.Velocity.y * 5);
            }
            rock.Update(water);
            if (rock.Position.y > Gdx.graphics.getHeight() + rockImage.getHeight()) {
                rock = null;
            }
        }

        cam.update();
        batch.setProjectionMatrix(cam.combined);

        checkKeyboard();

        water.Update();

        Gdx.gl.glClearColor(0, 0, 0, 0.5f);
        Gdx.gl.glClear(GL20.GL_COLOR_BUFFER_BIT);
        Gdx.gl.glBlendFunc(GL20.GL_SRC_ALPHA, GL20.GL_ONE_MINUS_SRC_ALPHA);

        batch.begin();
        batch.draw(backgroundImage, 0, 0, backgroundImage.getWidth(), backgroundImage.getHeight(), 0, 0,
                backgroundImage.getWidth(), backgroundImage.getHeight(), false, true);

        if (rock != null) {
            rock.Draw(batch, rockImage);
        }

        font.draw(batch, "Click the left mouse button to drop a rock!\n Rock dropped @x:" + x, 50, 50);
        batch.end();

        water.DrawToRenderTargets();

        water.Draw();

        batch.begin();
        batch.draw(background, 0, 200, background.getWidth(), background.getHeight(), 0, 0,
                background.getWidth(), background.getHeight(), false, true);
        batch.draw(background, 128 + 65, 200, background.getWidth(), background.getHeight(), 0, 0,
                background.getWidth(), background.getHeight(), false, true);
        batch.draw(background, 128 + 130, 200, background.getWidth(), background.getHeight(), 0, 0,
                background.getWidth(), background.getHeight(), false, true);
        batch.end();
    }

    private void checkKeyboard() {
        if (Gdx.input.isButtonPressed((Input.Buttons.LEFT)) || Gdx.input.isKeyPressed((Input.Keys.C))) {
            x = Gdx.input.getX();
            rock = new Rock(new Vector2(Gdx.input.getX(), Gdx.input.getY()), new Vector2(0.3f, 0.3f));
        }
    }

    @Override
    public void dispose() {
        batch.dispose();
        particleImage.dispose();
        backgroundImage.dispose();
        rockImage.dispose();
    }
}
