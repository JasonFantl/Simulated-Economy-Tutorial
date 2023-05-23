
class Community {

  FloatDict populations;


  Community(float foxCount, float rabbitCount, float grassCount) {
    populations = new FloatDict();
    populations.set("humans", 0);
    populations.set("foxes", foxCount);
    populations.set("rabbits", rabbitCount);
    populations.set("grass", grassCount);
  }

  void step() {
    
    if (economy != null) {
    populations.set("humans", economy.merchants.size());
    }
    int substeps = 400;
    for (int i = 0; i < substeps; i++) {
    FloatDict newPopulations = new FloatDict();

    for (String name : community.populations.keys()) {
      float newPopulationCount= populations.get(name) + changeInPopulation(name)*timeDelta/float(substeps);
      //if (newPopulationCount < 0.0) {
      //  newPopulationCount = 0.0;
      //}
      newPopulations.set(name, newPopulationCount);
    }
    populations = newPopulations;
    }
  }

  float changeInPopulation(String name) {
    if (name == "foxes") {
      return changeInfoxes();
    } else if (name == "rabbits") {
      return changeInRabbits();
    } else if (name == "grass") {
      return changeInGrass();
    }

    return 0.0;
  }

  float changeInfoxes() {
    // using the Lotka-Volterra model
    float c_FR = 0.005;
    float d_F = 0.1;
    float G = populations.get("grass");
    float R = populations.get("rabbits");
    float F = populations.get("foxes");
    
    //float populationDelta = c_FR*F*R - d_F*F;
    float z = 0.1;
    float y = 14;
    float populationDelta = z+F*(G-y);

    return populationDelta;
  }

  float changeInRabbits() {
    float c_RG = 0.001;
    float c_FR = 0.003;
    float G = populations.get("grass");
    float R = populations.get("rabbits");
    float F = populations.get("foxes");
    
    //float populationDelta = c_RG*R*G - c_FR*F*R;
    float z = 0.1;
    float populationDelta = G+z*R;

    return populationDelta;
  }

  float changeInGrass() {
    // using the logistic growth model
    float r_G = 1.5;
    float c_RG = 0.02;
    float K = 400;
    float G = populations.get("grass");
    float R = populations.get("rabbits");
    float F = populations.get("foxes");

    //float populationDelta = r_G*G*(1-G/K) - c_RG*G*R;
    float populationDelta = -R -F;

    return populationDelta;
  }
}
