class App extends React.Component {
  render() {
    if (this.loggedIn) {
      return (<LoggedIn />);
    } else {
      return (<LoggedIn />);
    }
  }
}

class Home extends React.Component {
  render() {
    return (
      <div className="container">
        <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
          <h1>Darezilla</h1>
          <p>A load of Dares XD</p>
          <p>Sign in to get access </p>
          <a onClick={this.authenticate} className="btn btn-primary btn-login btn-block">Sign In</a>
        </div>
      </div>
    )
  }
}

  class LoggedIn extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        dares: []
      };
      this.serverRequest = this.serverRequest.bind(this);
    }
    serverRequest() {
      $.get("http://localhost:3000/api/darez", res => {
        this.setState({
          dares: res
        });
      });
    }
    componentDidMount() {
      this.serverRequest();
    }
    render() {
      return (
        <div className="container">
          <div className="col-lg-12">
            <br />
            <span className="pull-right"><a onClick={this.logout}>Log out</a></span>
            <h2>Darezilla</h2>
            <p>Let's feed you with some fucken dares!!!</p>
            <div className="row">
              {this.state.dares.map(function(dare, i){
                return (<Dare key={i} dare={dare} />);
              })}
            </div>
          </div>
        </div>
      )
    }
  }

  class Dare extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        liked: "",
        dares: []
      };
      this.like = this.like.bind(this);
      this.serverRequest = this.serverRequest.bind(this);
    }
  
    like() {
      let dare = this.props.dare;
      this.serverRequest(dare);
    }
    serverRequest(dare) {
      $.post(
        "http://localhost:3000/api/darez/like/" + dare.id,
        { like: 1 },
        res => {
          console.log("res... ", res);
          this.setState({ liked: "Liked!", dares: res });
          this.props.dares = res;
        }
      );
    }
    
    render() {
      return (
        <div className="col-xs-4">
          <div className="panel panel-default">
            <div className="panel-heading">{this.props.dare.title}  </div>
            <div className="panel-body">
              {this.props.dare.dare}
            </div>
            <div className="panel-footer">
              {this.props.dare.likes} Likes &nbsp;
              <a onClick={this.like} className="btn btn-default">
                <span className="glyphicon glyphicon-thumbs-up"></span>
              </a>
            </div>
          </div>
        </div>
      )
    }
  }

ReactDOM.render(<App />, document.getElementById('app'));
