import React, { Component } from 'react';
import axios from 'axios';

export default class FilesUploadComponent extends Component {

    constructor(props) {
        super(props);

        this.onFileChange = this.onFileChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);

        this.state = {
            error: null,
            isReady: false,
            fileCollection: '',
            optimizationData: [],        
        }
    }

    onFileChange(e) {
        this.setState({ fileCollection: e.target.files })
    }

    onSubmit(e) {
        e.preventDefault()

        var formData = new FormData();
        for (const key of Object.keys(this.state.fileCollection)) {
            formData.append('file', this.state.fileCollection[key])
        }
        this.setState({isReady: false})
        axios.post("/api/v1/upload", formData, {})
        .then((response) => {
            this.setState({error: null, optimizationData: response.data, isReady: true})
        })
        .catch((err)=>{    
            this.setState({error: err})       
        })
    }

    render() {
        return (
            <div className="container">
                <div className="row" style={ {margin: '20px'}}>
                    <form onSubmit={this.onSubmit} >
                        <div className="form-group" style={ {padding: '10px'}}>
                            <input type="file" name="fileCollection" onChange={this.onFileChange} multiple />
                        </div>    
                        <div className="form-group" style={ {padding: '10px'}}>
                                <button className="btn btn-primary" type="submit">Start optimization</button>
                        </div> 
                        <div style={ {padding: '10px'}}>
                            { this.showResult() }
                        </div>                   
                    </form>
                </div>
            </div>
        )
    }

    showError() {
        return (
            <div>
                <div>Error: {this.state.error.message}</div>
                <div>Script response: {this.state.error.response.data.text}</div>
            </div>            
            ) 
    }

    showGenericError() {
        return (
            <div>
                <div>Error: {this.state.error.message}</div>
            </div>
        )
    }

    showScriptResult() {
        return (
            <div>
                <div>
                    <ul>
                        <li>Filename: {this.state.optimizationData.filename}</li>
                    </ul>
                    <ul>
                        <li>Location: {this.state.optimizationData.location}</li>
                    </ul>
                    <ul>
                        <li>ETag: {this.state.optimizationData.etag}</li>
                    </ul>
                </div>
            </div>
            )
    }

    showResult() {
        if (this.state.error) {
            if (this.state.error.response.data.text) {
                return this.showError()                 
            } else {
                return this.showGenericError()
            }
        } else if (!this.state.isReady) {
            if (this.state.optimizationData.length === 0) 
                return <div></div>;
            else
                return <div>Loading...</div>;
        } else {
            return this.showScriptResult()
        }
    }
}
