import Phaser from "phaser";
import MainScene from "./scenes/MainScene";
import MenuScene from "./scenes/MenuScene";

const SHARED_CONFIG = {
  width: 800,
  height: 600,
};

const config = {
  type: Phaser.AUTO,
  ...SHARED_CONFIG,
  physics: {
    default: "arcade",
    // arcade: {
    //   debug: true,
    // },
  },
  scene: [new MenuScene(SHARED_CONFIG), new MainScene(SHARED_CONFIG)],
};

new Phaser.Game(config);
