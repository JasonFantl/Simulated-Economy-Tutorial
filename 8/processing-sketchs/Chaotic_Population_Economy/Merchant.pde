import processing.data.JSONObject;
import java.util.HashMap;

class Merchant {
  float Money;
  String BuysSells;
  int CarryingCapacity;
  int Owned;
  HashMap<String, HashMap<String, Float>> ExpectedPrices;

  Merchant(String data) {
    JSONObject json = JSONObject.parse(data);
    this.Money = json.getFloat("Money");
    this.BuysSells = json.getString("BuysSells");
    this.CarryingCapacity = json.getInt("CarryingCapacity");
    this.Owned = json.getInt("Owned");
    this.ExpectedPrices = new HashMap<String, HashMap<String, Float>>();
    
    JSONObject jsonExpectedPrices = json.getJSONObject("ExpectedPrices");
    for (Object keyObject : jsonExpectedPrices.keys()) {
      String key = keyObject.toString();
      JSONObject jsonPrices = jsonExpectedPrices.getJSONObject(key);
      HashMap<String, Float> prices = new HashMap<String, Float>();
      for (Object subKeyObject : jsonPrices.keys()) {
        String subKey = subKeyObject.toString();
        prices.put(subKey, jsonPrices.getFloat(subKey));
      }
      this.ExpectedPrices.put(key, prices);
    }
  }

  JSONObject toJson() {
    JSONObject json = new JSONObject();
    json.setFloat("Money", this.Money);
    json.setString("BuysSells", this.BuysSells);
    json.setInt("CarryingCapacity", this.CarryingCapacity);
    json.setInt("Owned", this.Owned);
    
    JSONObject jsonExpectedPrices = new JSONObject();
    for (String key : this.ExpectedPrices.keySet()) {
      JSONObject jsonPrices = new JSONObject();
      HashMap<String, Float> prices = this.ExpectedPrices.get(key);
      for (String subKey : prices.keySet()) {
        jsonPrices.setFloat(subKey, prices.get(subKey));
      }
      jsonExpectedPrices.setJSONObject(key, jsonPrices);
    }
    json.setJSONObject("ExpectedPrices", jsonExpectedPrices);
    
    return json;
  }
}
