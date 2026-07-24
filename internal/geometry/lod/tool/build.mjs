import fs from "node:fs";
import mapshaper from "mapshaper";

const [input, output, resolution, weighting] = process.argv.slice(2);
if (!input || !output || !resolution || !weighting) {
  throw new Error("usage: node build.mjs INPUT OUTPUT RESOLUTION WEIGHTING");
}
const source = JSON.parse(fs.readFileSync(input, "utf8"));
if (source.type !== "FeatureCollection") throw new Error("input must be a FeatureCollection");
const features = [];
for (const feature of source.features) {
  const geometryID = feature.properties?.geometry_id;
  if (!geometryID) throw new Error("feature is missing geometry_id");
  const single = JSON.stringify({type: "FeatureCollection", features: [feature]});
  const command = `-i input.geojson -simplify resolution=${resolution} weighting=${weighting} planar keep-shapes -clean -o output.geojson format=geojson precision=0.000001`;
  const files = await mapshaper.applyCommands(command, {"input.geojson": single});
  const result = JSON.parse(String(files["output.geojson"]));
  if (result.type !== "FeatureCollection" || result.features.length !== 1) {
    throw new Error(`${geometryID}: expected one output feature, got ${result.features?.length ?? "invalid"}`);
  }
  result.features[0].properties = {geometry_id: geometryID};
  features.push(result.features[0]);
}
fs.writeFileSync(output, `${JSON.stringify({type: "FeatureCollection", features})}\n`);
