package restful

import (
	"io"
	"net/http/httptest"
	"testing"
)

func issue560RobotFile(req *Request, resp *Response) {
	io.WriteString(resp, "User-agent: *\nDisallow: /\n")
}

func issue560Hello(req *Request, resp *Response) {
	io.WriteString(resp, "hello "+req.PathParameter("name"))
}

// Test_detectWebServiceRegexServiceMustMatchToScore checks that a WebService
// with a regex path does not shadow a WebService without a path when its regex
// does not match the request token (issue #560).
func Test_detectWebServiceRegexServiceMustMatchToScore(t *testing.T) {
	router := CurlyRouter{}

	// new(WebService) is lazily initialized to Path("/") by Container.Add.
	defaultWS := new(WebService).Path("/")
	regexWS := new(WebService).Path("/{:hello}")
	wss := []*WebService{defaultWS, regexWS}

	// "robot.txt" does not match the regex of regexWS, so the plain defaultWS wins.
	if selected := router.detectWebService(tokenizePath("/robot.txt"), wss); selected != defaultWS {
		t.Fatalf("expected defaultWS, got %v", selected.rootPath)
	}

	// "hello" matches the regex of regexWS, so regexWS wins.
	if selected := router.detectWebService(tokenizePath("/hello/juan"), wss); selected != regexWS {
		t.Fatalf("expected regexWS, got %v", selected.rootPath)
	}
}

// TestCurlyRouterDefaultServiceStillServedWithRegexService is the end-to-end
// reproduction of issue #560: routes of a WebService without a path were never
// selected once a WebService with a regex path was registered, because the regex
// WebService scored higher even when its regex did not match.
func TestCurlyRouterDefaultServiceStillServedWithRegexService(t *testing.T) {
	defaultWS := new(WebService)
	defaultWS.Route(defaultWS.GET("/{:robot.txt}").To(issue560RobotFile))

	helloWS := new(WebService).Path("/{:hello}")
	helloWS.Route(helloWS.GET("/{name:*}").To(issue560Hello))

	c := NewContainer()
	c.Add(defaultWS)
	c.Add(helloWS)

	// The robot.txt route of the default WebService must be selected.
	robotRequest := httptest.NewRequest("GET", "/robot.txt", nil)
	robotResp := httptest.NewRecorder()
	c.ServeHTTP(robotResp, robotRequest)
	if robotResp.Code != 200 || robotResp.Body.String() != "User-agent: *\nDisallow: /\n" {
		t.Fatalf("expected 200 with the robot.txt content, got %d %q", robotResp.Code, robotResp.Body.String())
	}

	// The regex WebService must still handle its own requests.
	helloRequest := httptest.NewRequest("GET", "/hello/Juan", nil)
	helloResp := httptest.NewRecorder()
	c.ServeHTTP(helloResp, helloRequest)
	if helloResp.Code != 200 || helloResp.Body.String() != "hello Juan" {
		t.Fatalf("expected 200 with hello Juan, got %d %q", helloResp.Code, helloResp.Body.String())
	}
}
