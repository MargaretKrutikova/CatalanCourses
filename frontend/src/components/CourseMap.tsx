import { useState } from "react";
import { CourseMarker } from "./CourseMarker";
import { Map, MapCameraChangedEvent } from "@vis.gl/react-google-maps";
import courseData from "../data/course_list.json";

const courses = courseData.data.filter((c) => parseInt(c.horaInici.split(":")[0]) > 16);
const home = {
  lat: parseFloat(process.env.REACT_APP_HOME_LAT || ""),
  lng: parseFloat(process.env.REACT_APP_HOME_LNG || ""),
};

const courseMarkers = courses.map((course) => ({
  key: course.codiPlain,
  location: {
    lat: course.latitud,
    lng: course.longitud,
  },
}));

export const CourseMap: React.FunctionComponent<{}> = () => {
  const [activeMarkerIndex, setActiveMarkerIndex] = useState<number | null>(null);

  const setActiveMarker = (index: number | null) => {
    setActiveMarkerIndex(index);
  };

  return (
    <Map
      defaultZoom={14}
      mapId={"6a73725fb27f7adc"}
      defaultCenter={home}
      onCameraChanged={(ev: MapCameraChangedEvent) =>
        console.log("camera changed:", ev.detail.center, "zoom:", ev.detail.zoom)
      }
    >
      {courseMarkers.map((marker, index) => (
        <CourseMarker
          key={marker.key}
          handleInfoClose={() => setActiveMarker(null)}
          marker={marker}
          setActive={() => setActiveMarker(index)}
          isActive={activeMarkerIndex === index}
          course={courses[index]}
        />
      ))}
    </Map>
  );
};
