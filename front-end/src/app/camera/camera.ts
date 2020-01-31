/**
 * Camera object, which has a name and link to its feed.
 */
export class Camera {
  name: string;
  link: string;

  constructor(jsonData) {
    this.name = jsonData.name;
    this.link = jsonData.link;
  }
}
