import "./App.css";
import { APIProvider } from "@vis.gl/react-google-maps";
import { CourseMap } from "./components/CourseMap";

const App = () => {
  return (
    <div className="App">
      <header className="App-header"></header>
      <APIProvider
        apiKey={process.env.REACT_APP_GOOGLE_API_KEY!}
        onLoad={() => console.log("Maps API has loaded.")}
      >
        <CourseMap />
      </APIProvider>
    </div>
  );
};

export default App;
