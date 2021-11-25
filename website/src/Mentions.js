import React from "react";
import './App.css';
import CountUp from 'react-countup';


/*API REQUEST */
class Mentions extends React.Component {
   
    // Constructor 
    constructor(props) {
        super(props);

        this.state = {
            mentions: 0,
            prevMentions: 0,
            DataisLoaded: false
        };

    }

    pullData() {
        fetch("https://api.axolobot.ml/v1/info")
            .then((res) => res.json())
            .then((json) => {
                this.setState({
                    prevMentions: this.state.mentions,
                    mentions: json["mentions"],
                    DataisLoaded: true
                });
                console.log(json)
            }
        )
    }

    // ComponentDidMount is used to execute the code 
    componentDidMount() {
        this.interval = setInterval(() => this.pullData(), 5000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    mentionsStyle = {
        color: 'white',
        fontFamily: 'consolas',
        backgroundColor: '#f69fa7',
        padding: '5px',
        margin: '2px',
        display: 'inline-block',
        borderRadius: '10px',
    }

    loadingStyle = {
        display: 'inline-block',
    }




    render() {
        const { DataisLoaded, mentions, prevMentions } = this.state;
        if (!DataisLoaded) {
            return (<div style={this.mentionsStyle}>
                ---
            </div>);
        }
        return (<div style={this.mentionsStyle}>
                    <CountUp start={prevMentions} end={mentions} delay={0} duration="3" useEasing="true" />
                </div>);
    }
  }

  export default Mentions;