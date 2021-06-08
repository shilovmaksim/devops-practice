import React, { Component } from 'react';
import '../node_modules/bootstrap/dist/css/bootstrap.min.css';
import FilesUploadComponent from './components/files-upload-component';
import {
  BrowserRouter as Router,
  Switch,
  Route
} from "react-router-dom";

class App extends Component {
  render() {
    return(
    <Router>
        {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
        <Switch>
          <Route path="/health">
             <div>ok</div>
          </Route>
          <Route path="/">
            <div className="App">
              <FilesUploadComponent />
            </div>
          </Route>
        </Switch>
    </Router>
    )
  }
}

export default App;
