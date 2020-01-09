export class Camera {
  name: string;
  link: string;

  constructor(jsonData) {
    this.name = jsonData.name;
    this.link = jsonData.link;
  }
}
