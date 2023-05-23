
class Community {

  FloatDict populations;


  Community(float foxCount, float rabbitCount, float grassCount) {
    populations = new FloatDict();
    populations.set("foxs", foxCount);
    populations.set("rabbits", rabbitCount);
    populations.set("grass", grassCount);
  }

  void step() {
    
    int substeps = 400;
    for (int i = 0; i < substeps; i++) {
    FloatDict newPopulations = new FloatDict();

    for (String name : community.populations.keys()) {
      float newPopulationCount= populations.get(name) + changeInPopulation(name)*timeDelta/float(substeps);
      if (newPopulationCount < 0.0) {
        newPopulationCount = 0.0;
      }
      newPopulations.set(name, newPopulationCount);
    }
    populations = newPopulations;
    }
  }

  float changeInPopulation(String name) {
    if (name == "foxs") {
      return changeInFoxs();
    } else if (name == "rabbits") {
      return changeInRabbits();
    } else if (name == "grass") {
      return changeInGrass();
    }

    return 0.0;
  }

  float changeInFoxs() {
    // using the Lotka-Volterra model
    float c_FR = 0.003;
    float d_F = 0.1;
    float R = populations.get("rabbits");
    float F = populations.get("foxs");
    
    float populationDelta = c_FR*F*R - d_F*F;

    return populationDelta;
  }

  float changeInRabbits() {
    float c_RG = 0.001;
    float c_FR = 0.003;
    float G = populations.get("grass");
    float R = populations.get("rabbits");
    float F = populations.get("foxs");
    
    float populationDelta = c_RG*R*G - c_FR*F*R;

    return populationDelta;
  }

  float changeInGrass() {
    // using the logistic growth model
    float r_G = 1.5;
    float c_RG = 0.02;
    float K = 1000;
    float G = populations.get("grass");
    float R = populations.get("rabbits");

    float populationDelta = r_G*G*(1-G/K) - c_RG*G*R;
    //float populationDelta = r_G*G*(1-G/K) - c_RG*G*R;

    return populationDelta;
  }
}
