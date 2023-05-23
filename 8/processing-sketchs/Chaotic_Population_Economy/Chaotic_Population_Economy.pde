import gifAnimation.*;
GifMaker gifExport;

import grafica.*;

// create economy connection
EconomyInterface economy;

// Initialize the community and graph
Community community;
GPlot plot, plot2;

// Variables to control the scrolling of the graph
float scrollSpeed = 1;
float timeCounter = 0;
int maxPoints = 160;

float timeDelta = 0.05;
float tickerFrac = 4.0;
float timePassed = 0.0;

void setup() {
  size(1000, 500);
  frameRate(60);
  gifExport = new GifMaker(this, "export.gif");
  gifExport.setRepeat(0); // make it an "endless" animation
  gifExport.setTransparent(0, 0, 0); // make black the transparent color. every black pixel in the animation will be transparent
  
  community = new Community(20, -10, 10.1);

  // Set up the plot
  plot = new GPlot(this);
  plot.setPos(0, 0);
  plot.setDim(400, 400);
  plot.setXLim(0, maxPoints);
  plot.setAxesOffset(0);
  plot.getYAxis().setDrawTickLabels(true);
  plot.getYAxis().setAxisLabelText("Populations");
  plot.getXAxis().setAxisLabelText("Time");
  plot.setXLim(0.0, maxPoints*timeDelta*tickerFrac);
  plot.setYLim(-80, 80);

  for (String name : community.populations.keys()) {
    plot.addLayer(name, new GPointsArray(maxPoints));
    plot.getLayer(name).setLineWidth(3);

    if (name == "foxes") {
      plot.getLayer(name).setLineColor(color(250, 100, 25));
    } else if (name == "rabbits") {
      plot.getLayer(name).setLineColor(color(20, 10, 250));
    } else if (name == "grass") {
      plot.getLayer(name).setLineColor(color(20, 200, 25));
    } else {
      plot.getLayer(name).setLineColor(color(200, 200, 25));
    }
  }

  plot2 = new GPlot(this);
  plot2.setPos(480, 0);
  plot2.setDim(400, 400);
  plot2.getTitle().setText("Phase space");
  plot2.getXAxis().getAxisLabel().setText("Rabbits");
  plot2.getYAxis().getAxisLabel().setText("Grass");
  plot2.setXLim(-50, 50);
  plot2.setYLim(-50, 50);
}

void draw() {
  background(255);

  if (economy != null) {
  economy.update();
  }

  boolean movedXAxis= false;
  for (String name : community.populations.keys()) {
    plot.getLayer(name).addPoint(new GPoint(timePassed*tickerFrac, community.populations.get(name)));

    // Remove the first point if the number of points exceeds maxPoints
    if (plot.getLayer(name).getPointsRef().getNPoints() > maxPoints) {
      plot.getLayer(name).removePoint(0);

      if (!movedXAxis) {
        float newLowerLim = plot.getXLim()[0] + timeDelta*tickerFrac;
        float newUpperLim = plot.getXLim()[1] + timeDelta*tickerFrac;
        plot.setXLim(newLowerLim, newUpperLim);
        movedXAxis = true;
      }
    }
  }

  // Update the community
  community.step();

  // Draw the first plot
  plot.beginDraw();
  plot.drawBackground();
  plot.drawBox();
  plot.drawXAxis();
  plot.drawYAxis();
  plot.drawTopAxis();
  plot.drawRightAxis();
  plot.drawTitle();
  plot.drawLines();
  plot.drawLabels();
  plot.endDraw();



  plot2.addPoint(community.populations.get("rabbits"), community.populations.get("grass"));
  if (plot2.getPointsRef().getNPoints() > maxPoints) {
    plot2.removePoint(0);
  }

  // Draw the second plot
  plot2.beginDraw();
  plot2.drawBackground();
  plot2.drawBox();
  plot2.drawXAxis();
  plot2.drawYAxis();
  plot2.drawTitle();
  plot2.drawGridLines(GPlot.BOTH);
  plot2.drawLines();
  plot2.drawPoint(plot2.getPointsRef().getLastPoint());
  plot2.endDraw();

  timePassed += timeDelta;

  //gifExport.addFrame();
}

void keyPressed() {
  //gifExport.finish();
  //println("gif saved");
    economy = new EconomyInterface(this, "localhost", 55555);

}
