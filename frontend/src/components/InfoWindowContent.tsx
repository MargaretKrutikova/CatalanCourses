import detailedCourseInfo from "../data/complete_course_info.json";
import { CourseType } from "../Types";

const detailedCourseInfoMap = detailedCourseInfo.reduce<Record<string, (typeof detailedCourseInfo)[number]>>(
  (acc, info) => ({ ...acc, ...{ [info.code]: info } }),
  {}
);

export type InfoWindowContentProps = {
  course: CourseType;
};

export const InfoWindowContent = ({ course }: InfoWindowContentProps) => {
  const detailedInfo = detailedCourseInfoMap[course.codiPlain] || {};

  return (
    <>
      <div>{course.codiPlain}</div>
      <div>{course.dies}</div>
      <div>
        {course.horaInici} - {course.horaFi}
      </div>
      <div>
        Places total {course.placesLliures}, places left {detailedInfo.placesLeft}
      </div>
      <div>{detailedInfo.address}</div>
      <div>{detailedInfo.metros}</div>
      <div>Deadline: {detailedInfo.registrationDeadline}</div>
    </>
  );
};
