import Phaser from "phaser";

class MenuScene extends Phaser.Scene {
	constructor(config) {
		super("MenuScene");
		this.config = config;
		this.menu = [
			{ scene: "MainScene", text: "Play" },
			{ scene: null, text: "Exit" },
		];
	}

	preload() {
		this.load.image("bg", "assets/sky.png");
	}

	create() {
		this.add.image(0, 0, "bg").setOrigin(0, 0);
		let margin = 0;

		this.add
			.text(this.config.width / 2, this.config.height / 4, "Berakout", {
				fontSize: "100px",
				fill: "#99692e",
			})
			.setOrigin(0.5);

		this.add
			.text(
				this.config.width / 2,
				this.config.height - 40,
				"create by aji mustofa @pepega90",
				{
					fontSize: "20px",
					fill: "#333",
				}
			)
			.setOrigin(0.5);

		this.menu.forEach(m => {
			let textPos = [
				this.config.width / 2,
				this.config.height / 2 + margin,
			];
			let textMenu = this.add
				.text(...textPos, m.text, {
					fontSize: "32px",
					fill: "#fff",
				})
				.setOrigin(0.5);
			margin += 60;

			textMenu.setInteractive();

			textMenu.on("pointerover", () => {
				textMenu.setStyle({ fill: "#ff0" });
			});

			textMenu.on("pointerout", () => {
				textMenu.setStyle({ fill: "#fff" });
			});

			textMenu.on("pointerup", () => {
				m.scene && this.scene.start(m.scene);

				if (m.text == "Exit") this.game.destroy(true);
			});
		});
	}
}

export default MenuScene;
