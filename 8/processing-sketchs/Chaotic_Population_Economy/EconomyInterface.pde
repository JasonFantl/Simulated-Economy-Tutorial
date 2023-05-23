import processing.net.*;


class EconomyInterface {
  private Client client;
  boolean connected;
  String name;

  ArrayList<Merchant> merchants;
  float averagebedPrice;

  EconomyInterface(PApplet parent, String serverAddress, int serverPort) {
    client = new Client(parent, serverAddress, serverPort);
    merchants = new ArrayList<Merchant>();

    // send city name
    name = "JAVAVILLE";
    sendMessage(name);
  }

  void sendMessage(String message) {
    client.write(message + "\n");
  }

  String receiveMessage() {
    return client.readString().trim();
  }

  void sendMerchantData(JSONObject merchantData) {
    String jsonString = merchantData.toString();
    jsonString = jsonString.replaceAll("\\s+", ""); // Remove whitespace characters
    sendMessage(jsonString);
  }

  void update() {
    // wait until the city name has been recieved
    if (!connected) {
      if (client.available() > 0) {
        println("Connected to " + client.readString().trim());
        connected = true;
      }
      return;
    }

    // recieve merchants
    if (client.available() > 0) {
      Merchant m = new Merchant(receiveMessage());

      merchants.add(m);
      //println(m.toJson());
    }

    ArrayList<Merchant> merchantsToRemove = new ArrayList<>();
    for (Merchant m : merchants) {
      // leave the city if reached carrying capacity
      if (m.Owned >= m.CarryingCapacity) {
        merchantsToRemove.add(m);
        m.ExpectedPrices.get("fur").put(name, 2.0);

        sendMerchantData(m.toJson());
        continue;
      }

      // get a fur
      if (random(1.0) < (community.populations.get("foxes")-1)/10) {
        m.Owned+=5;
        community.populations.set("foxes", community.populations.get("foxes")-1);
      }
    }
    merchants.removeAll(merchantsToRemove);
  }

  void close() {
    client.stop();
  }
}
