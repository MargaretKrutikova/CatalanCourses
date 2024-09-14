export type CourseType = {
  codi: string;
  codiPlain: string;
  nom: string;
  idNivell: number;
  modalitat: string;
  idModalitat: number;
  municipi: string;
  dies: string;
  horaInici: string;
  horaFi: string;
  matriculacio: string;
  placesLliures: number;
  longitud: number;
};

export type CourseMarkerType = {
  key: string;
  location: {
    lat: number;
    lng: number;
  };
};
