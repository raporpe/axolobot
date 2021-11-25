import React from "react";
import logo from './logo.png';
import twitterlogo from './twitter-logo.png'
import './App.css';
import phonephoto from './background-phone.png';
import explanation from './explanation.jpg';
import team from './team.png';
import Mentions from './Mentions.js'
import Confetti from "./Confetti";

function App() {
  return (

    <div>
    <meta charSet="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
    <meta name="description" content />
    <meta name="author" content />
    <title>AxoloBot ✨</title>
    <link rel="icon" type="image/x-icon" href="assets/favicon.ico" />
    {/* Bootstrap icons*/}
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.5.0/font/bootstrap-icons.css" rel="stylesheet" />
    {/* Google fonts*/}
    <link rel="preconnect" href="https://fonts.gstatic.com" />
    <link href="https://fonts.googleapis.com/css2?family=Newsreader:ital,wght@0,600;1,600&display=swap" rel="stylesheet" />
    <link href="https://fonts.googleapis.com/css2?family=Mulish:ital,wght@0,300;0,500;0,600;0,700;1,300;1,500;1,600;1,700&display=swap" rel="stylesheet" />
    <link href="https://fonts.googleapis.com/css2?family=Kanit:ital,wght@0,400;1,400&display=swap" rel="stylesheet" />
    {/* Core theme CSS (includes Bootstrap)*/}
    <link href="css/styles.css" rel="stylesheet" />
    {/* Navigation*/}
    <nav className="navbar navbar-expand-lg navbar-light fixed-top shadow-sm" id="mainNav">
      <div className="container px-5">
        <a className="navbar-brand fw-bold" href="#page-top"><img src={logo} alt="..." style={{height: '3rem'}} /> AxoloBot</a>
        
        <div className="collapse navbar-collapse" id="navbarResponsive">
          <ul className="navbar-nav ms-auto me-4 my-3 my-lg-0">
            <li className="nav-item"><a className="nav-link me-lg-3" href="#features">How To Use</a></li>
            <li className="nav-item"><a className="nav-link me-lg-3" href="#team">The Team</a></li>
          </ul>
         
        </div>
        <a href="https://twitter.com/axolobot"> 
          <button className="btn btn-primary rounded-pill px-3 mb-2 mb-lg-0" data-bs-toggle="modal" data-bs-target="#feedbackModal">
            <span className="d-flex align-items-center">
              <i className="bi-lightning-fill me-2" />
              <span className="small">Try Now</span>
            </span>
          </button>
          </a>
      </div>
    </nav>
    {/* Mashead header*/}
    <header className="masthead">
      <div className="container px-5">
        <div className="row gx-5 align-items-center">
          <div className="col-lg-6">
            {/* Mashead text and app badges*/}
            <div className="mb-5 mb-lg-0 text-center text-lg-start">
              <h1 className="display-1 lh-1 mb-3">Discover what they think.</h1>
              <p className="lead fw-normal text-muted mb-5">AxoloBot is your perfect companion, it uses powerful AI techniques to know a tweet's reaction in few seconds. <br/>How? It gathers and analyse responses!</p>
              <a href="https://github.com/raporpe/axolobot"> 
                <button className="btn btn-primary rounded-pill px-3 mb-2 mb-lg-0 black" data-bs-toggle="modal" data-bs-target="#feedbackModal">
                  <span className="d-flex align-items-center">
                    <i className="bi-github me-2" />
                    <span className="small"> Watch on Github</span>
                  </span>
                </button>
              </a>
              {/*<div className="d-flex flex-column flex-lg-row align-items-center">
                <a className="me-lg-3 mb-4 mb-lg-0" href="#!"><img className="app-badge" src="assets/img/google-play-badge.svg" alt="..." /></a>
                <a href="#!"><img className="app-badge" src="assets/img/app-store-badge.svg" alt="..." /></a>
              </div>*/}
            </div>
          </div>
          <div className="col-lg-6">
            {/* Masthead device mockup feature*/}
            <div className="masthead-device-mockup">
              <svg className="circle" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                <defs>
                  <linearGradient id="circleGradient" gradientTransform="rotate(45)">
                    <stop className="gradient-start-color" offset="0%" />
                    <stop className="gradient-end-color" offset="100%" />
                  </linearGradient>
                </defs>
                <circle cx={50} cy={50} r={50} /></svg><svg className="shape-1 d-none d-sm-block" viewBox="0 0 240.83 240.83" xmlns="http://www.w3.org/2000/svg">
                <rect x="-32.54" y="78.39" width="305.92" height="84.05" rx="42.03" transform="translate(120.42 -49.88) rotate(45)" />
                <rect x="-32.54" y="78.39" width="305.92" height="84.05" rx="42.03" transform="translate(-49.88 120.42) rotate(-45)" /></svg><svg className="shape-2 d-none d-sm-block" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><circle cx={50} cy={50} r={50} /></svg>
              <div className="device-wrapper">
                <div className="device" data-device="iPhoneX" data-orientation="portrait" data-color="black">
                  <div className="screen bg-black">
                    <img source src={phonephoto} alt="..." style={{maxWidth: '100%', height: '100%'}}/>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
    {/* Quote/testimonial aside*/}
    <aside className="text-center bg-gradient-primary-to-secondary">
      <div className="container px-5">
        <div className="row gx-5 justify-content-center">
          <div className="col-xl-8">
            <div className="h2 fs-1 text-white mb-4">Axolobot has been used <Mentions></Mentions> times</div>
            <img src={twitterlogo} alt="..." style={{height: '2.2rem'}} />
          </div>
        </div>
      </div>
    </aside>
    {/* App features section*/}
    <section id="features">
      <div className="container px-5">
        <div className="row gx-5 align-items-center">
          <div className="col-lg-8 order-lg-1 mb-5 mb-lg-0">
            <div className="container-fluid px-5">
              <div className="row gx-5">
                <div className="col-md-6 mb-5">
                  {/* Feature item*/}
                  <div className="text-center">
                    <i className="bi-search icon-feature text-gradient d-block mb-3" />
                    <h3 className="font-alt">1. Look for the tweet</h3>
                    <p className="text-muted mb-0">The tweet must be not older than 7 days, written in English or Spanish with at least 1 reply.</p>
                  </div>
                </div>
                <div className="col-md-6 mb-5">
                  {/* Feature item*/}
                  <div className="text-center">
                    <i className="bi-chat-square-quote-fill icon-feature text-gradient d-block mb-3" />
                    <h3 className="font-alt">2. Mention @axolobot</h3>
                    <p className="text-muted mb-0">Write a comment to the tweet mentioning the coolest bot of all.</p>
                  </div>
                </div>
              </div>
              <div className="row">
                <div className="col-md-6 mb-5 mb-md-0">
                  {/* Feature item*/}
                  <div className="text-center">
                    <i className="bi-sun-fill icon-feature text-gradient d-block mb-3" />
                    <h3 className="font-alt">3. Watch the analysis</h3>
                    <p className="text-muted mb-0">AxoloBot will answer you with the average generated feeling.</p>
                  </div>
                </div>
                <div className="col-md-6">
                  {/* Feature item*/}
                  <div className="text-center">
                    <i className="bi-heart-fill icon-feature text-gradient d-block mb-3" />
                    <h3 className="font-alt">4. Share the knowledge</h3>
                    <p className="text-muted mb-0">AxoloBot can help others,<br/> do not keep it as a secret!</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="col-lg-4 order-lg-0">
            {/* Features section device mockup*/}
            <div className="features-device-mockup">
              <svg className="circle" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                <defs>
                  <linearGradient id="circleGradient" gradientTransform="rotate(45)">
                    <stop className="gradient-start-color" offset="0%" />
                    <stop className="gradient-end-color" offset="100%" />
                  </linearGradient>
                </defs>
                <circle cx={50} cy={50} r={50} /></svg><svg className="shape-1 d-none d-sm-block" viewBox="0 0 240.83 240.83" xmlns="http://www.w3.org/2000/svg">
                <rect x="-32.54" y="78.39" width="305.92" height="84.05" rx="42.03" transform="translate(120.42 -49.88) rotate(45)" />
                <rect x="-32.54" y="78.39" width="305.92" height="84.05" rx="42.03" transform="translate(-49.88 120.42) rotate(-45)" /></svg><svg className="shape-2 d-none d-sm-block" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><circle cx={50} cy={50} r={50} /></svg>
              <div className="device-wrapper">
                <div className="device" data-device="iPhoneX" data-orientation="portrait" data-color="black">
                  <div className="screen bg-black">
                    {/* PUT CONTENTS HERE:*/}
                    {/* * * This can be a video, image, or just about anything else.*/}
                    {/* * * Set the max width of your media to 100% and the height to*/}
                    {/* * * 100% like the demo example below.*/}
                    <img source src={explanation} alt="..." style={{maxWidth: '100%', height: '100%'}}/>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    {/* Basic features section*/}
    <section className="bg-light">
      <div className="container px-5"  id="team">
        <div className="row gx-5 align-items-center justify-content-center justify-content-lg-between">
          <div className="col-12 col-lg-5">
            <h2 className="display-4 lh-1 mb-4">It is great. <br/> Its team, too.</h2>
            <p className="lead fw-normal text-muted mb-5 mb-lg-0">Three great students willing to use AI in a helpful way for everyone ended up creating AxoloBot. Check them out! <br/><br/>
          </p>

          <div>
            <div>
          <a href="https://www.linkedin.com/in/raul-portugues/"> 
              <button className="btn btn-primary rounded-pill px-3 mb-2 mb-lg-0" data-bs-toggle="modal" data-bs-target="#feedbackModal">
              <span className="d-flex align-items-center">
              <i className="bi-linkedin me-2" />
              <span className="small">Raúl's Linkedin</span>
            </span>
          </button>
          </a>
              </div>

              <div>
          <a href="https://www.linkedin.com/in/diegogarciaperez/"> 
              <button className="btn btn-primary rounded-pill px-3 mb-2 mb-lg-0" data-bs-toggle="modal" data-bs-target="#feedbackModal">
              <span className="d-flex align-items-center">
              <i className="bi-linkedin me-2" />
              <span className="small">Diego's Linkedin</span>
            </span>
          </button>
          </a>
          </div>

          <div>
          <a href="https://www.linkedin.com/in/jorge-lizcano/"> 
              <button className="btn btn-primary rounded-pill px-3 mb-2 mb-lg-0" data-bs-toggle="modal" data-bs-target="#feedbackModal">
              <span className="d-flex align-items-center">
              <i className="bi-linkedin me-2" />
              <span className="small">Jorge's Linkedin</span>
            </span>
          </button>
          </a>
          </div>
          </div>
          </div>
          
          <div className="col-sm-8 col-md-6">
            <div className="px-5 px-sm-0"><img className="img-fluid" src={team} alt="..." /></div>
          </div>
        </div>
      </div>
    </section>
    {/* App badge section*/}
    <section className="bg-gradient-primary-to-secondary">
      <div className="container px-5">
        <h2 className="text-center text-white font-alt mb-4">Thanks for your support! <br/> AxoloBot will be waiting for your tasty tweets around here.</h2>
        <p className="text-center text-white font-alt mb-4 credits">Credits to Enrique Monroy for the awesome logo &#9834; </p>
      </div>
    </section>
    {/* Footer*/}
    <footer className="bg-black text-center py-5">
      <div className="container px-5">
        <div className="text-white-50 small">
          <div className="mb-2">© AxoloBot 2021</div>
        </div>
      </div>
    </footer>
    {/* Feedback Modal*/}
    <div className="modal fade" id="feedbackModal" tabIndex={-1} aria-labelledby="feedbackModalLabel" aria-hidden="true">
      <div className="modal-dialog modal-dialog-centered">
        <div className="modal-content">
            <div className="modal-header bg-gradient-primary-to-secondary p-4">
              <h5 className="modal-title font-alt text-white" id="feedbackModalLabel">Try Now</h5>
              <button className="btn-close btn-close-white" type="button" data-bs-dismiss="modal" aria-label="Close" />
            </div>
          <div className="modal-body border-0 p-4">
            {/* * * * * * * * * * * * * * * **/}
            {/* * * SB Forms Contact Form * **/}
            {/* * * * * * * * * * * * * * * **/}
            {/* This form is pre-integrated with SB Forms.*/}
            {/* To make this form functional, sign up at*/}
            {/* https://startbootstrap.com/solution/contact-forms*/}
            {/* to get an API token!*/}
            <form id="contactForm" data-sb-form-api-token="API_TOKEN">
              {/* Name input*/}
              <div className="form-floating mb-3">
                <input className="form-control" id="name" type="text" placeholder="Enter your name..." data-sb-validations="required" />
                <label htmlFor="name">Full name</label>
                <div className="invalid-feedback" data-sb-feedback="name:required">A name is required.</div>
              </div>
              {/* Email address input*/}
              <div className="form-floating mb-3">
                <input className="form-control" id="email" type="email" placeholder="name@example.com" data-sb-validations="required,email" />
                <label htmlFor="email">Email address</label>
                <div className="invalid-feedback" data-sb-feedback="email:required">An email is required.</div>
                <div className="invalid-feedback" data-sb-feedback="email:email">Email is not valid.</div>
              </div>
              {/* Phone number input*/}
              <div className="form-floating mb-3">
                <input className="form-control" id="phone" type="tel" placeholder="(123) 456-7890" data-sb-validations="required" />
                <label htmlFor="phone">Phone number</label>
                <div className="invalid-feedback" data-sb-feedback="phone:required">A phone number is required.</div>
              </div>
              {/* Message input*/}
              <div className="form-floating mb-3">
                <textarea className="form-control" id="message" type="text" placeholder="Enter your message here..." style={{height: '10rem'}} data-sb-validations="required" defaultValue={""} />
                <label htmlFor="message">Message</label>
                <div className="invalid-feedback" data-sb-feedback="message:required">A message is required.</div>
              </div>
              {/* Submit success message*/}
              {/**/}
              {/* This is what your users will see when the form*/}
              {/* has successfully submitted*/}
              <div className="d-none" id="submitSuccessMessage">
                <div className="text-center mb-3">
                  <div className="fw-bolder">Form submission successful!</div>
                  To activate this form, sign up at
                  <br />
                  <a href="https://startbootstrap.com/solution/contact-forms">https://startbootstrap.com/solution/contact-forms</a>
                </div>
              </div>
              {/* Submit error message*/}
              {/**/}
              {/* This is what your users will see when there is*/}
              {/* an error submitting the form*/}
              <div className="d-none" id="submitErrorMessage"><div className="text-center text-danger mb-3">Error sending message!</div></div>
              {/* Submit Button*/}
              <div className="d-grid"><button className="btn btn-primary rounded-pill btn-lg disabled" id="submitButton" type="submit">Submit</button></div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
  );
}


export default App;



