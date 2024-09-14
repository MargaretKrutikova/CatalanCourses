import { AdvancedMarker, InfoWindow, Pin, useAdvancedMarkerRef } from "@vis.gl/react-google-maps";
import { InfoWindowContent } from "./InfoWindowContent";
import { CourseMarkerType, CourseType } from "../Types";

export type CourseMarkerProps = {
  marker: CourseMarkerType;
  course: CourseType;
  setActive: () => void;
  handleInfoClose: () => void;
  isActive: Boolean;
};

export const CourseMarker = ({ marker, course, setActive, isActive, handleInfoClose }: CourseMarkerProps) => {
  const [markerRef, markerElement] = useAdvancedMarkerRef();

  return (
    <>
      <AdvancedMarker position={marker.location} ref={markerRef} onClick={setActive}>
        <Pin background={"#FBBC04"} glyphColor={"#000"} borderColor={"#000"} />
        {isActive ? (
          <InfoWindow anchor={markerElement} onClose={handleInfoClose} onCloseClick={handleInfoClose}>
            <InfoWindowContent course={course} />
          </InfoWindow>
        ) : null}
      </AdvancedMarker>
    </>
  );
};
