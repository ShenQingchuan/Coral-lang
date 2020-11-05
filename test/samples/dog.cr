import "animal.cr"

class Dog : Animal {
  var color String;

  fn Dog(name String, age int, color String) {
    super(name, age);
    this.color = color;
  }

  public fn greet() {
    printf("Hi, I'm a %s dog named %s!", this.color, this.name);
  }
}

var d1 = new Dog("John", 6, "#bbb");
d1.greet();
