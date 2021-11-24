import React from "react";
import './App.css';


/*API REQUEST */
class Request extends React.Component {
   
    // Constructor 
    constructor(props) {
        super(props);

        this.state = {
            mentions: 0,
            DataisLoaded: false
        };
    }

    pullData() {
        fetch("https://api.axolobot.ml/v1/info")
            .then((res) => res.json())
            .then((json) => {
                this.setState({
                    mentions: json["mentions"],
                    DataisLoaded: true
                });
                console.log(json)
            }
        )
    }

    // ComponentDidMount is used to execute the code 
    componentDidMount() {
        this.interval = setInterval(() => this.pullData(), 1000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    render() {
        const { DataisLoaded, mentions } = this.state;
        if (!DataisLoaded) {
        return (<div><h1>Loading tweets...</h1></div>);
        }
        return (<div className = "App">AxoloBot has been used {mentions} times!</div>);
    }
  }

  export default Request;