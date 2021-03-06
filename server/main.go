/******************************************************************************
 *
 *  Description :
 *
 *  Setup & initialization.
 *
 *****************************************************************************/

package main

//go:generate protoc --proto_path=../pbx --go_out=plugins=grpc:../pbx ../pbx/model.proto

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

  "context"
  "github.com/tinode/chat/pbx"

	// For stripping comments from JSON config
	jcr "github.com/DisposaBoy/JsonConfigReader"

	gh "github.com/gorilla/handlers"

	// Authenticators
	"github.com/tinode/chat/server/auth"
	_ "github.com/tinode/chat/server/auth/anon"
	_ "github.com/tinode/chat/server/auth/basic"
	_ "github.com/tinode/chat/server/auth/rest"
	_ "github.com/tinode/chat/server/auth/token"

	// Database backends
	_ "github.com/tinode/chat/server/db/mysql"
	_ "github.com/tinode/chat/server/db/rethinkdb"

	// Push notifications
	"github.com/tinode/chat/server/push"
	_ "github.com/tinode/chat/server/push/fcm"
	_ "github.com/tinode/chat/server/push/stdout"

	"github.com/tinode/chat/server/store"

	// Credential validators
	_ "github.com/tinode/chat/server/validate/email"
	_ "github.com/tinode/chat/server/validate/tel"
	"google.golang.org/grpc"

	// File upload handlers
	_ "github.com/tinode/chat/server/media/fs"
	_ "github.com/tinode/chat/server/media/s3"


	"encoding/gob"
	"fmt"
	"os/signal"
	"syscall"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

var isConnEmpty bool = true
var isConnEmpty2 bool = true

type waHandler struct {
	c *whatsapp.Conn
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *waHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

//Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (*waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Printf("%v %v %v %v\n\t%v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.Info.QuotedMessageID, message.Text)

    sub := &pbx.ClientSub{}
    sub.Topic = "usrNoJ5tCr-JCM"
    // sub.Content = []byte("hihi")
    msgSub := &pbx.ClientMsg_Sub{sub}
    clientMessage3 := &pbx.ClientMsg{Message: msgSub}
    err3 := globals.stream.Send(clientMessage3)
    if err3 != nil {
      log.Fatal("error sending message ", err3)
    }

    note := &pbx.ClientNote{}
    note.Topic = "usrNoJ5tCr-JCM"
    note.What = 2
    msgNote := &pbx.ClientMsg_Note{note}
    clientMessage5 := &pbx.ClientMsg{Message: msgNote}
    err5 := globals.stream.Send(clientMessage5)
    if err5 != nil {
      log.Fatal("error sending message ", err5)
    }


    pub := &pbx.ClientPub{}
    pub.Topic = "usrNoJ5tCr-JCM"
    pub.NoEcho = true

    bytes, err4 := json.Marshal(message.Text)
    if err4 != nil {
      log.Fatal("error sending message ", err4)
    }

    pub.Content = bytes
    msgPub := &pbx.ClientMsg_Pub{pub}
    clientMessage2 := &pbx.ClientMsg{Message: msgPub}
    err2 := globals.stream.Send(clientMessage2)
    if err2 != nil {
      log.Fatal("error sending message ", err2)
    }



	return


  // sess, count := globals.sessionStore.NewSession(ws, "")
  // log.Println("ws: session started", sess.sid, count)


//	var msg ClientComMessage
//	s.lastAction = types.TimeNow()
//	msg.timestamp = s.lastAction
//
//	msg.from = s.uid.UserId()
//	msg.authLvl = int(s.authLvl)

  // log.Println("HandleTextMessage: session started", globals.sessionStore)

  var err error
  // var stream pbx.Node_MessageLoopClient
  if isConnEmpty {
    // isConnEmpty = false

  	fmt.Println("globals.conn == nil")

    globals.conn, err = grpc.Dial("localhost:6061", grpc.WithInsecure())
    if err != nil {
      isConnEmpty = true
      log.Fatal("Error dialing", err)
    }
    defer globals.conn.Close()

    client := pbx.NewNodeClient(globals.conn)

    // ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
    defer cancel()

    stream, err := client.MessageLoop(ctx)
    // response, err := client.MessageLoop(context.Background())

    if err != nil {
      log.Fatal("Error calling", err)
    }



    hi := &pbx.ClientHi{}
    hi.Id = "1"
    hi.UserAgent = "Golang_Spider_Bot/3.0"
    hi.Ver = "0.15"
    hi.Lang = "EN"

    msgHi := &pbx.ClientMsg_Hi{hi}
    clientMessage := &pbx.ClientMsg{Message: msgHi}
    err = stream.Send(clientMessage)

    if err != nil {
      log.Fatal("error sending message ", err)
    }

    login := &pbx.ClientLogin{}
    // login.Id = "xena"
    login.Scheme = "basic"
    login.Secret = []byte("xena:xena123")
    clMsg := &pbx.ClientMsg_Login{login}
    clientMessage = &pbx.ClientMsg{Message: clMsg}
    err = stream.Send(clientMessage)

    if err != nil {
      log.Fatal("error sending message ", err)
    }


    pub := &pbx.ClientPub{}
    pub.Topic = "usrNoJ5tCr-JCM"
    pub.Content = []byte("")
    msgPub := &pbx.ClientMsg_Pub{pub}
    clientMessage2 := &pbx.ClientMsg{Message: msgPub}
    err2 := stream.Send(clientMessage2)
    if err2 != nil {
      log.Fatal("error sending message ", err2)
    }


//     serverMsg, err := stream.Recv()
//     if err != nil {
//       log.Fatal(err)
//     }
//     log.Println(serverMsg)
// 
//     serverMsg, err = stream.Recv()
//     if err != nil {
//       log.Fatal(err)
//     }
//     log.Println(serverMsg)

    waitc := make(chan struct{})
    go func() {
      for {
        in, err := stream.Recv()
        if err == io.EOF {
          // read done.
          close(waitc)
          return
        }
        if err != nil {
          log.Fatalf("Failed to receive a note : %v", err)
        }
        log.Printf("Got message %s", in)
      }
    }()
    // for _, note := range notes {
    //   if err := stream.Send(note); err != nil {
    //     log.Fatalf("Failed to send a note: %v", err)
    //   }
    // }
    // stream.CloseSend()
    <-waitc


//   } else if isConnEmpty2 {
//     isConnEmpty2 = false
// 
//     pub := &pbx.ClientPub{}
//     pub.Topic = "usrNoJ5tCr-JCM"
//     pub.Content = []byte(message.Text)
//     pubMsg := &pbx.ClientMsg_Pub{pub}
//     clientMessage := &pbx.ClientMsg{Message: pubMsg}
//     err = stream.Send(clientMessage)
//     if err != nil {
//       log.Fatal("error sending message ", err)
//     }


  }


	fmt.Printf("%v %v %v %v\n\t%v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.Info.QuotedMessageID, message.Text)
}

/*//Example for media handling. Video, Audio, Document are also possible in the same way
func (*waHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	data, err := message.Download()
	if err != nil {
		return
	}
	filename := fmt.Sprintf("%v/%v.%v", os.TempDir(), message.Info.Id, strings.Split(message.Type, "/")[1])
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	log.Printf("%v %v\n\timage reveived, saved at:%v\n", message.Info.Timestamp, message.Info.RemoteJid, filename)
}*/

const (
	// idleSessionTimeout defines duration of being idle before terminating a session.
	idleSessionTimeout = time.Second * 55
	// idleTopicTimeout defines now long to keep topic alive after the last session detached.
	idleTopicTimeout = time.Second * 5

	// currentVersion is the current API/protocol version
	currentVersion = "0.15"
	// minSupportedVersion is the minimum supported API version
	minSupportedVersion = "0.15"

	// defaultMaxMessageSize is the default maximum message size
	defaultMaxMessageSize = 1 << 19 // 512K

	// defaultMaxSubscriberCount is the default maximum number of group topic subscribers.
	// Also set in adapter.
	defaultMaxSubscriberCount = 256

	// defaultMaxTagCount is the default maximum number of indexable tags
	defaultMaxTagCount = 16

	// minTagLength is the shortest acceptable length of a tag in runes. Shorter tags are discarded.
	minTagLength = 4
	// maxTagLength is the maximum length of a tag in runes. Longer tags are trimmed.
	maxTagLength = 96

	// Delay before updating a User Agent
	uaTimerDelay = time.Second * 5

	// maxDeleteCount is the maximum allowed number of messages to delete in one call.
	defaultMaxDeleteCount = 1024

	// Mount point where static content is served, http://host-name/<defaultStaticMount>
	defaultStaticMount = "/"

	// Local path to static content
	defaultStaticPath = "static"
)

// Build version number defined by the compiler:
// 		-ldflags "-X main.buildstamp=value_to_assign_to_buildstamp"
// Reported to clients in response to {hi} message.
// For instance, to define the buildstamp as a timestamp of when the server was built add a
// flag to compiler command line:
// 		-ldflags "-X main.buildstamp=`date -u '+%Y%m%dT%H:%M:%SZ'`"
var buildstamp = "undef"

// CredValidator holds additional config params for a credential validator.
type credValidator struct {
	// AuthLevel(s) which require this validator.
	requiredAuthLvl []auth.Level
	addToTags       bool
}

var globals struct {
  wac          *whatsapp.Conn

  conn         *grpc.ClientConn
  stream       pbx.Node_MessageLoopClient


	hub          *Hub
	sessionStore *SessionStore
	cluster      *Cluster
	grpcServer   *grpc.Server
	plugins      []Plugin
	statsUpdate  chan *varUpdate

	// Credential validators.
	validators map[string]credValidator
	// Validators required for each auth level.
	authValidators map[auth.Level][]string

	apiKeySalt []byte
	// Tag namespaces (prefixes) which are immutable to the client.
	immutableTagNS map[string]bool
	// Tag namespaces which are immutable on User and partially mutable on Topic:
	// user can only mutate tags he owns.
	maskedTagNS map[string]bool

	// Add Strict-Transport-Security to headers, the value signifies age.
	// Empty string "" turns it off
	tlsStrictMaxAge string
	// Listen for connections on this address:port and redirect them to HTTPS port.
	tlsRedirectHTTP string
	// Maximum message size allowed from peer.
	maxMessageSize int64
	// Maximum number of group topic subscribers.
	maxSubscriberCount int
	// Maximum number of indexable tags.
	maxTagCount int

	// Maximum allowed upload size.
	maxFileUploadSize int64
}

type validatorConfig struct {
	// TRUE or FALSE to set
	AddToTags bool `json:"add_to_tags"`
	//  Authentication level which triggers this validator: "auth", "anon"... or ""
	Required []string `json:"required"`
	// Validator params passed to validator unchanged.
	Config json.RawMessage `json:"config"`
}

type mediaConfig struct {
	// The name of the handler to use for file uploads.
	UseHandler string `json:"use_handler"`
	// Maximum allowed size of an uploaded file
	MaxFileUploadSize int64 `json:"max_size"`
	// Garbage collection timeout
	GcPeriod int `json:"gc_period"`
	// Number of entries to delete in one pass
	GcBlockSize int `json:"gc_block_size"`
	// Individual handler config params to pass to handlers unchanged.
	Handlers map[string]json.RawMessage `json:"handlers"`
}

// Contentx of the configuration file
type configType struct {
	// Default HTTP(S) address:port to listen on for websocket and long polling clients. Either a
	// numeric or a canonical name, e.g. ":80" or ":https". Could include a host name, e.g.
	// "localhost:80".
	// Could be blank: if TLS is not configured, will use ":80", otherwise ":443".
	// Can be overridden from the command line, see option --listen.
	Listen string `json:"listen"`
	// Cache-Control value for static content.
	CacheControl int `json:"cache_control"`
	// Address:port to listen for gRPC clients. If blank gRPC support will not be initialized.
	// Could be overridden from the command line with --grpc_listen.
	GrpcListen string `json:"grpc_listen"`
	// URL path for mounting the directory with static files.
	StaticMount string `json:"static_mount"`
	// Local path to static files. All files in this path are made accessible by HTTP.
	StaticData string `json:"static_data"`
	// Salt used in signing API keys
	APIKeySalt []byte `json:"api_key_salt"`
	// Maximum message size allowed from client. Intended to prevent malicious client from sending
	// very large files inband (does not affect out of band uploads).
	MaxMessageSize int `json:"max_message_size"`
	// Maximum number of group topic subscribers.
	MaxSubscriberCount int `json:"max_subscriber_count"`
	// Masked tags: tags immutable on User (mask), mutable on Topic only within the mask.
	MaskedTagNamespaces []string `json:"masked_tags"`
	// Maximum number of indexable tags
	MaxTagCount int `json:"max_tag_count"`
	// URL path for exposing runtime stats. Disabled if the path is blank.
	ExpvarPath string `json:"expvar"`

	// Configs for subsystems
	Cluster   json.RawMessage             `json:"cluster_config"`
	Plugin    json.RawMessage             `json:"plugins"`
	Store     json.RawMessage             `json:"store_config"`
	Push      json.RawMessage             `json:"push"`
	TLS       json.RawMessage             `json:"tls"`
	Auth      map[string]json.RawMessage  `json:"auth_config"`
	Validator map[string]*validatorConfig `json:"acc_validation"`
	Media     *mediaConfig                `json:"media"`
}

func main() {
	executable, _ := os.Executable()

	// All relative paths are resolved against the executable path, not against current working directory.
	// Absolute paths are left unchanged.
	rootpath, _ := filepath.Split(executable)

	log.Printf("Server v%s:%s:%s; db: '%s'; pid %d; %d process(es)",
		currentVersion, executable, buildstamp,
		store.GetAdapterName(), os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))

	var configfile = flag.String("config", "tinode.conf", "Path to config file.")
	// Path to static content.
	var staticPath = flag.String("static_data", defaultStaticPath, "Path to directory with static files to be served.")
	var listenOn = flag.String("listen", "", "Override address and port to listen on for HTTP(S) clients.")
	var listenGrpc = flag.String("grpc_listen", "", "Override address and port to listen on for gRPC clients.")
	var tlsEnabled = flag.Bool("tls_enabled", false, "Override config value for enabling TLS.")
	var clusterSelf = flag.String("cluster_self", "", "Override the name of the current cluster node.")
	var expvarPath = flag.String("expvar", "", "Override the path where runtime stats are exposed.")
	var pprofFile = flag.String("pprof", "", "File name to save profiling info to. Disabled if not set.")
	flag.Parse()

	*configfile = toAbsolutePath(rootpath, *configfile)
	log.Printf("Using config from '%s'", *configfile)

	var config configType
	if file, err := os.Open(*configfile); err != nil {
		log.Fatal("Failed to read config file: ", err)
	} else if err = json.NewDecoder(jcr.New(file)).Decode(&config); err != nil {
		log.Fatal("Failed to parse config file: ", err)
	}

	if *listenOn != "" {
		config.Listen = *listenOn
	}

	// Initialize cluster and receive calculated workerId.
	// Cluster won't be started here yet.
	workerId := clusterInit(config.Cluster, clusterSelf)

	if *pprofFile != "" {
		*pprofFile = toAbsolutePath(rootpath, *pprofFile)

		cpuf, err := os.Create(*pprofFile + ".cpu")
		if err != nil {
			log.Fatal("Failed to create CPU pprof file: ", err)
		}
		defer cpuf.Close()

		memf, err := os.Create(*pprofFile + ".mem")
		if err != nil {
			log.Fatal("Failed to create Mem pprof file: ", err)
		}
		defer memf.Close()

		pprof.StartCPUProfile(cpuf)
		defer pprof.StopCPUProfile()
		defer pprof.WriteHeapProfile(memf)

		log.Printf("Profiling info saved to '%s.(cpu|mem)'", *pprofFile)
	}

	err := store.Open(workerId, string(config.Store))
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer func() {
		store.Close()
		log.Println("Closed database connection(s)")
		log.Println("All done, good bye")
	}()

	// API key signing secret
	globals.apiKeySalt = config.APIKeySalt

	err = store.InitAuthLogicalNames(config.Auth["logical_names"])
	if err != nil {
		log.Fatal(err)
	}

	authNames := store.GetAuthNames()
	for _, name := range authNames {
		if authhdl := store.GetLogicalAuthHandler(name); authhdl == nil {
			log.Fatalln("Unknown authenticator", name)
		} else if jsconf := config.Auth[name]; jsconf != nil {
			if err := authhdl.Init(string(jsconf), name); err != nil {
				log.Fatalln("Failed to init auth scheme", name+":", err)
			}
		}
	}

	// List of tag namespaces for user discovery which cannot be changed directly
	// by the client, e.g. 'email' or 'tel'.
	globals.immutableTagNS = make(map[string]bool)

	// Process validators.
	for name, vconf := range config.Validator {
		// Check if validator is restrictive. If so, add validator name to the list of restricted tags.
		// The namespace can be restricted even if the validator is disabled.
		if vconf.AddToTags {
			if strings.Contains(name, ":") {
				log.Fatal("acc_validation names should not contain character ':'")
			}
			globals.immutableTagNS[name] = true
		}

		if len(vconf.Required) == 0 {
			// Skip disabled validator.
			continue
		}

		var reqLevels []auth.Level
		for _, req := range vconf.Required {
			lvl := auth.ParseAuthLevel(req)
			if lvl == auth.LevelNone {
				if req != "" {
					log.Fatalf("Invalid required AuthLevel '%s' in validator '%s'", req, name)
				}
				// Skip empty string
				continue
			}
			reqLevels = append(reqLevels, lvl)
			if globals.authValidators == nil {
				globals.authValidators = make(map[auth.Level][]string)
			}
			globals.authValidators[lvl] = append(globals.authValidators[lvl], name)
		}

		if len(reqLevels) == 0 {
			// Ignore validator with empty levels.
			continue
		}

		if val := store.GetValidator(name); val == nil {
			log.Fatal("Config provided for an unknown validator '" + name + "'")
		} else if err = val.Init(string(vconf.Config)); err != nil {
			log.Fatal("Failed to init validator '"+name+"': ", err)
		}
		if globals.validators == nil {
			globals.validators = make(map[string]credValidator)
		}
		globals.validators[name] = credValidator{
			requiredAuthLvl: reqLevels,
			addToTags:       vconf.AddToTags}
	}

	// Partially restricted tag namespaces
	globals.maskedTagNS = make(map[string]bool, len(config.MaskedTagNamespaces))
	for _, tag := range config.MaskedTagNamespaces {
		if strings.Contains(tag, ":") {
			log.Fatal("masked_tags namespaces should not contain character ':'")
		}
		globals.maskedTagNS[tag] = true
	}

	var tags []string
	for tag := range globals.immutableTagNS {
		tags = append(tags, "'"+tag+"'")
	}
	if len(tags) > 0 {
		log.Println("Restricted tags:", tags)
	}
	tags = nil
	for tag := range globals.maskedTagNS {
		tags = append(tags, "'"+tag+"'")
	}
	if len(tags) > 0 {
		log.Println("Masked tags:", tags)
	}

	// Maximum message size
	globals.maxMessageSize = int64(config.MaxMessageSize)
	if globals.maxMessageSize <= 0 {
		globals.maxMessageSize = defaultMaxMessageSize
	}
	// Maximum number of group topic subscribers
	globals.maxSubscriberCount = config.MaxSubscriberCount
	if globals.maxSubscriberCount <= 1 {
		globals.maxSubscriberCount = defaultMaxSubscriberCount
	}
	// Maximum number of indexable tags per user or topics
	globals.maxTagCount = config.MaxTagCount
	if globals.maxTagCount <= 0 {
		globals.maxTagCount = defaultMaxTagCount
	}

	if config.Media != nil {
		if config.Media.UseHandler == "" {
			config.Media = nil
		} else {
			globals.maxFileUploadSize = config.Media.MaxFileUploadSize
			if config.Media.Handlers != nil {
				var conf string
				if params := config.Media.Handlers[config.Media.UseHandler]; params != nil {
					conf = string(params)
				}
				if err = store.UseMediaHandler(config.Media.UseHandler, conf); err != nil {
					log.Fatalf("Failed to init media handler '%s': %s", config.Media.UseHandler, err)
				}
			}
			if config.Media.GcPeriod > 0 && config.Media.GcBlockSize > 0 {
				stopFilesGc := largeFileRunGarbageCollection(time.Second*time.Duration(config.Media.GcPeriod),
					config.Media.GcBlockSize)
				defer func() {
					stopFilesGc <- true
					log.Println("Stopped files garbage collector")
				}()
			}
		}
	}

	err = push.Init(string(config.Push))
	if err != nil {
		log.Fatal("Failed to initialize push notifications:", err)
	}
	defer func() {
		push.Stop()
		log.Println("Stopped push notifications")
	}()

	// Keep inactive LP sessions for 15 seconds
	globals.sessionStore = NewSessionStore(idleSessionTimeout + 15*time.Second)
	// The hub (the main message router)
	globals.hub = newHub()

	// Start accepting cluster traffic.
	if globals.cluster != nil {
		globals.cluster.start()
	}

	tlsConfig, err := parseTLSConfig(*tlsEnabled, config.TLS)
	if err != nil {
		log.Fatalln(err)
	}

	// Intialize plugins
	pluginsInit(config.Plugin)

	// Set up gRPC server, if one is configured
	if *listenGrpc == "" {
		*listenGrpc = config.GrpcListen
	}
	if globals.grpcServer, err = serveGrpc(*listenGrpc, tlsConfig); err != nil {
		log.Fatal(err)
	}

	// Set up HTTP server. Must use non-default mux because of expvar.
	mux := http.NewServeMux()

	// Serve static content from the directory in -static_data flag if that's
	// available, otherwise assume '<path-to-executable>/static'. The content is served at
	// the path pointed by 'static_mount' in the config. If that is missing then it's
	// served at root '/'.
	var staticMountPoint string
	if *staticPath != "" && *staticPath != "-" {
		// Resolve path to static content.
		*staticPath = toAbsolutePath(rootpath, *staticPath)
		if _, err = os.Stat(*staticPath); os.IsNotExist(err) {
			log.Fatal("Static content directory is not found", *staticPath)
		}

		staticMountPoint = config.StaticMount
		if staticMountPoint == "" {
			staticMountPoint = defaultStaticMount
		} else {
			if !strings.HasPrefix(staticMountPoint, "/") {
				staticMountPoint = "/" + staticMountPoint
			}
			if !strings.HasSuffix(staticMountPoint, "/") {
				staticMountPoint += "/"
			}
		}
		mux.Handle(staticMountPoint,
			// Add optional Cache-Control header
			cacheControlHandler(config.CacheControl,
				// Optionally add Strict-Transport_security to the response
				hstsHandler(
					// Add gzip compression
					gh.CompressHandler(
						// And add custom formatter of errors.
						httpErrorHandler(
							// Remove mount point prefix
							http.StripPrefix(staticMountPoint,
								http.FileServer(http.Dir(*staticPath))))))))
		log.Printf("Serving static content from '%s' at '%s'", *staticPath, staticMountPoint)
	} else {
		log.Println("Static content is disabled")
	}

	// Handle websocket clients.
	mux.HandleFunc("/v0/channels", serveWebSocket)
	// Handle long polling clients. Enable compression.
	mux.Handle("/v0/channels/lp", gh.CompressHandler(http.HandlerFunc(serveLongPoll)))
	if config.Media != nil {
		// Handle uploads of large files.
		mux.Handle("/v0/file/u/", gh.CompressHandler(http.HandlerFunc(largeFileUpload)))
		// Serve large files.
		mux.Handle("/v0/file/s/", gh.CompressHandler(http.HandlerFunc(largeFileServe)))
		log.Println("Large media handling enabled", config.Media.UseHandler)
	}

	if staticMountPoint != "/" {
		// Serve json-formatted 404 for all other URLs
		mux.HandleFunc("/", serve404)
	}

	evpath := *expvarPath
	if evpath == "" {
		evpath = config.ExpvarPath
	}
	statsInit(mux, evpath)


	//create new WhatsApp connection
	globals.wac, err = whatsapp.NewConn(5 * time.Second)
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}

	//Add handler
	globals.wac.AddHandler(&waHandler{globals.wac})

	//login or restore
	if err := login(globals.wac); err != nil {
		log.Fatalf("error logging in: %v\n", err)
	}

	go func() {
  	log.Println("XXXXXXXXXXXXXXXXXXXXXX");

    var err error

    globals.conn, err = grpc.Dial("localhost:6061", grpc.WithInsecure())
    if err != nil {
      isConnEmpty = true
      log.Fatal("Error dialing", err)
    }
    defer globals.conn.Close()

    client := pbx.NewNodeClient(globals.conn)

    // ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
    defer cancel()

    globals.stream, err = client.MessageLoop(ctx)
    // response, err := client.MessageLoop(context.Background())

    if err != nil {
      log.Fatal("Error calling", err)
    }



    hi := &pbx.ClientHi{}
    hi.Id = "1"
    hi.UserAgent = "Golang_Spider_Bot/3.0"
    hi.Ver = "0.15"
    hi.Lang = "EN"

    msgHi := &pbx.ClientMsg_Hi{hi}
    clientMessage := &pbx.ClientMsg{Message: msgHi}
    err = globals.stream.Send(clientMessage)

    if err != nil {
      log.Fatal("error sending message ", err)
    }

    login := &pbx.ClientLogin{}
    // login.Id = "xena"
    login.Scheme = "basic"
    login.Secret = []byte("xena:xena123")
    clMsg := &pbx.ClientMsg_Login{login}
    clientMessage = &pbx.ClientMsg{Message: clMsg}
    err = globals.stream.Send(clientMessage)

    if err != nil {
      log.Fatal("error sending message ", err)
    }



//     serverMsg, err := globals.stream.Recv()
//     if err != nil {
//       log.Fatal(err)
//     }
//     log.Println(serverMsg)
// 
//     serverMsg, err = globals.stream.Recv()
//     if err != nil {
//       log.Fatal(err)
//     }
//     log.Println(serverMsg)

    waitc := make(chan struct{})
    go func() {
      for {
        in, err := globals.stream.Recv()
        if err == io.EOF {
          // read done.
          close(waitc)
          return
        }
        if err != nil {
          log.Fatalf("Failed to receive a note : %v", err)
        }
        log.Printf("Got message %s", in)
      }
    }()
    // for _, note := range notes {
    //   if err := globals.stream.Send(note); err != nil {
    //     log.Fatalf("Failed to send a note: %v", err)
    //   }
    // }
    // globals.stream.CloseSend()
    <-waitc

	}()


	if err = listenAndServe(config.Listen, mux, tlsConfig, signalHandler()); err != nil {
		log.Fatal(err)
	}




	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	//Disconnect safe
	fmt.Println("Shutting down now.")
	session, err := globals.wac.Disconnect()
	if err != nil {
		log.Fatalf("error disconnecting: %v\n", err)
	}
	if err := writeSession(session); err != nil {
		log.Fatalf("error saving session: %v", err)
	}



}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
