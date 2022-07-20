import Phaser from "phaser";

class MainScene extends Phaser.Scene {
	constructor(config) {
		super("MainScene");
		this.config = config;

		this.SPEED = 500;
		this.ballSpeed = { x: Phaser.Math.Between(-50, 50), y: -300 };

		this.ballOnPlayer = true;
		this.player;
		this.ball;
		this.listWall;
		this.lives = 3;
		this.liveText;
		this.gameOverText;
		this.restartText;
		this.introText;
		this.gameOver = false;
		this.completeText;

		// keyboard input
		this.pressLeft;
		this.pressRight;
		this.pressRestart;
	}

	preload() {
		this.load.image("bg", "assets/sky.png");
		this.load.image("player", "assets/paddleRed.png");
		this.load.image("ball", "assets/ballBlue.png");
		this.load.image("ballLive", "assets/ballBlue.png");
		this.load.image("wall", "assets/element_red_rectangle_glossy.png");
		this.load.image("wall2", "assets/element_yellow_rectangle_glossy.png");
		this.load.image("wall3", "assets/element_green_rectangle_glossy.png");
	}

	create() {
		this.add.image(0, 0, "bg").setOrigin(0, 0);
		this.add.image(20, this.config.height - 15, "ballLive").setOrigin(0, 1);
		this.liveText = this.add
			.text(50, this.config.height - 15, "x " + this.lives, {
				fontSize: "25px",
				fill: "#333",
			})
			.setOrigin(0, 1);
		console.log(this);
		this.completeText = this.add
			.text(
				this.config.width / 2,
				this.config.height / 4,
				"Thanks for playing!",
				{
					fontSize: "50px",
					fill: "#99692e",
				}
			)
			.setOrigin(0.5);
		this.completeText.visible = false;

		this.player = this.physics.add
			.sprite(this.config.width / 2, this.config.height - 50, "player")
			.setImmovable();
		this.player.body.collideWorldBounds = true;
		this.pressRight = this.input.keyboard.addKey(
			Phaser.Input.Keyboard.KeyCodes.RIGHT
		);
		this.pressLeft = this.input.keyboard.addKey(
			Phaser.Input.Keyboard.KeyCodes.LEFT
		);
		this.pressRestart = this.input.keyboard.addKey(
			Phaser.Input.Keyboard.KeyCodes.R
		);

		this.ball = this.physics.add.sprite(
			this.config.width / 2,
			this.player.y - 25,
			"ball"
		);

		this.ball.body.bounce.setTo(1, 1);

		this.listWall = this.physics.add.group();
		this.createWall();

		this.introText = this.add
			.text(
				this.config.width / 2,
				this.config.height / 2 + 100,
				'Tekan "SPACE" untuk start game',
				{ fontSize: 38, fill: "#99692e" }
			)
			.setOrigin(0.5);

		this.gameOverText = this.add
			.text(
				this.config.width / 2,
				this.config.height / 2 + 100,
				"Game Over",
				{ fontSize: 38, fill: "#fc0303" }
			)
			.setOrigin(0.5);
		this.restartText = this.add
			.text(
				this.config.width / 2,
				this.config.height / 2 + 150,
				'Tekan "R" untuk restart game',
				{ fontSize: 25, fill: "#fc0303" }
			)
			.setOrigin(0.5);

		this.gameOverText.visible = false;
		this.restartText.visible = false;
	}

	update(time, delta) {
		this.input.keyboard.on("keydown_SPACE", () => {
			this.introText.visible = false;
			if (!this.gameOver) {
				this.ball.body.velocity.y = this.ballSpeed.y;
				this.ball.body.velocity.x = this.ballSpeed.x;
				this.ballOnPlayer = false;
			}
		});

		if (this.ballOnPlayer) {
			this.ball.body.x = this.player.x - 10;
		}

		if (!this.gameOver) {
			if (this.pressRight.isDown)
				this.player.body.velocity.x = this.SPEED;
			else if (this.pressLeft.isDown)
				this.player.body.velocity.x = -this.SPEED;
			else this.player.body.velocity.x = 0;
		}

		if (this.pressRestart.isDown) {
			this.scene.restart();
			this.lives = 3;
			this.gameOver = false;
			this.gameOverText.visible = false;
			this.restartText.visible = false;
		}

		if (
			this.ball.body.x > this.config.width - this.ball.width ||
			this.ball.body.x <= 0
		)
			this.ball.body.velocity.x *= -1;
		if (this.ball.body.y < 0) this.ball.body.velocity.y *= -1;

		if (this.ball.body.y > this.config.height) {
			this.ballOnPlayer = true;
			this.ball.body.reset(this.player.x - 10, this.player.y - 25);
			this.lives -= 1;
			this.liveText.setText("x " + this.lives);
		}

		if (this.lives < 1) {
			this.gameOver = true;
			this.gameOverText.visible = true;
			this.restartText.visible = true;
		}

		if (this.listWall.getChildren().length < 1) {
			this.completeText.visible = true;
		}

		this.physics.add.collider(
			this.player,
			this.ball,
			this.ballHitPlayer,
			null,
			this
		);
		this.physics.add.collider(
			this.ball,
			this.listWall,
			this.ballHitWall,
			null,
			this
		);
	}

	createWall() {
		let wall;

		for (let y = 0; y < 6; y++) {
			for (let x = 0; x < 9; x++) {
				if (y < 2) {
					wall = this.listWall.create(
						100 + x * 70,
						50 + y * 52,
						"wall"
					);
				} else if (y < 4) {
					wall = this.listWall.create(
						100 + x * 70,
						50 + y * 52,
						"wall2"
					);
				} else if (y < 6) {
					wall = this.listWall.create(
						100 + x * 70,
						50 + y * 52,
						"wall3"
					);
				}
				wall.body.bounce.set(1);
				wall.body.immovable = true;
			}
		}
	}

	ballHitPlayer() {
		let diff = 0;

		if (this.ball.x < this.player.x) {
			// ketika bola berada di kiri player
			diff = this.player.x - this.ball.x;
			this.ball.body.velocity.x = -5 * diff;
		} else if (this.ball.x > this.player.x) {
			// ketika bola berada di kanan player
			diff = this.ball.x - this.player.x;
			this.ball.body.velocity.x = 5 * diff;
		} else {
			// ketika bola di tengah-tengah
			this.ball.body.velocity.x = 2 + Math.random() * 8;
		}
	}

	ballHitWall(ball, wall) {
		wall.destroy();
	}
}

export default MainScene;
