package dbus

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/cylonchau/gofirewallder/object"
)

var (
	dbusClient          *DbusClientSerivce
	remotelyBusLck      sync.Mutex
	PORT                = 55557
	DEFAULT_ZONE_TARGET = "{chain}_{zone}"
)

type DbusClientSerivce struct {
	Conn      *dbus.Conn
	defaultZone string
}

func NewDbusClientService(addr string) (*DbusClientSerivce, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	remotelyBusLck.Lock()
	defer remotelyBusLck.Unlock()
	if dbusClient != nil && dbusClient.Conn.Connected() {
		return dbusClient, nil
	}
	conn, err := dbus.Connect("tcp:host="+host+",port="+port, dbus.WithAuth(dbus.AuthAnonymous()))
	if err != nil {
		return nil, err
	}

	obj := conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.INTERFACE_GETDEFAULTZONE, dbus.FlagNoAutoStart)

	if call.Err != nil {
		return nil, call.Err
	}

	return &DbusClientSerivce{
		conn,
		call.Body[0].(string),
	}, err
}

// @title         Reload
// @description   temporary Add rich language rule into zone.
// @auth      	  author           2021-10-05
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) generatePath(zone, interface_path string) (path dbus.ObjectPath, err error) {

	zoneid := c.getZoneId(zone)
	if zoneid < 0 {
		return "", errors.New("invalid zone.")
	}
	return dbus.ObjectPath(fmt.Sprintf("%s/%d", interface_path, zoneid)), nil
}

func (c *DbusClientSerivce) GetDefaultZone() string {
	return c.defaultZone
}

// @title         GetZoneSettings
// @description   Return runtime settings of given zone.
// @auth      	  author           2021-09-26
// @return        zones            []string       "Return array of names (s) of predefined zones known to current runtime environment."
// @return        error            error          ""
func (c *DbusClientSerivce) GetZones() (zones []string, err error) {
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_GETZONES, dbus.FlagNoAutoStart)

	if call.Err != nil {
		return nil, call.Err
	}
	return call.Body[0].([]string), nil
}

// @title         GetZoneSettings
// @description   Return runtime settings of given zone.
// @auth      	  author           2021-09-26
// @param         zone		       string         "zone name."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) getZoneId(zone string) int {
	var (
		zoneArray []string
		err       error
	)
	if zoneArray, err = c.GetZones(); err != nil {
		return -1
	}
	index := sort.SearchStrings(zoneArray, zone)

	if index < len(zoneArray) && zoneArray[index] == zone {
		return index
	} else {
		return -1
	}
}

// @title         GetZoneSettings
// @description   Return runtime settings of given zone.
// @auth      	  author           2021-09-26
// @param         zone		       string         "zone name."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) GetZoneSettings(zone string) (err error) {
	if err = c.checkZoneName(zone); err != nil {
		return err
	}

	obj := c.Conn.Object(object.INTERFACE, object.PATH)

	call := obj.Call(object.INTERFACE_GETZONESETTINGS, dbus.FlagNoAutoStart, zone)
	if call.Err != nil {
		return call.Err
	}

	return
}

// @title         AddZone
// @description   Add zone with given settings into permanent configuration.
// @auth      	  author           2021-09-27
// @param         name		       string         "Is an optional start and end tag and is used to give a more readable name."
// @return        error            error          "Possible errors: NAME_CONFLICT, INVALID_NAME, INVALID_TYPE"

func (c *DbusClientSerivce) AddZone(name string) (err error) {
	if err = c.checkZoneName(name); err != nil {
		return err
	}

	obj := c.Conn.Object(object.INTERFACE, object.CONFIG_PATH)
	zoneSettings := Settings{}

	zoneSettings.Targe = "default"
	zoneSettings.Short = name

	call := obj.Call(object.CONFIG_ADDZONE, dbus.FlagNoAutoStart, name, zoneSettings)

	if call.Err != nil {
		return call.Err
	}
	return
}

// @title         GetZoneOfInterface
// @description   temporary add a firewalld port
// @auth      	  author           2021-09-27
// @param         iface    		   string         "e.g. eth0, iface is device name."
// @return        zoneName         string         "Return name (s) of zone the interface is bound to or empty string.."
func (c *DbusClientSerivce) GetZoneOfInterface(iface string) string {
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_GETZONEOFINTERFACE, dbus.FlagNoAutoStart, iface)
	return call.Body[0].(string)
}

/************************************************** port area ***********************************************************/

// @title         addPort
// @description   temporary add a firewalld port
// @auth      	  author           2021-09-29
// @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
// @param         zone    		   string         "e.g. public|dmz.. The empty string is usage default zone, is currently firewalld defualt zone"
// @param         timeout    	   int	          "Timeout, 0 is the permanent effect of the currently service startup state."
// @return        zoneName         string         "Returns name of zone to which the protocol was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_PORT, MISSING_PROTOCOL, INVALID_PROTOCOL, ALREADY_ENABLED, INVALID_COMMAND."

func (c *DbusClientSerivce) AddPort(port, zone string, timeout int) (list string, err error) {

	if err = checkPort(port); err != nil {
		return "", err
	}

	if zone == "" {
		zone = c.GetDefaultZone()
	}

	port, protocol := splitPortProtocol(port)

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_ADDPORT, dbus.FlagNoAutoStart, zone, port, protocol, timeout)

	if call.Err != nil {
		return "", call.Err
	}
	return call.Body[0].(string), nil
}

// @title         PermanentAddPort
// @description   Permanently add port & procotol to list of ports of zone.
// @auth      	  author           2021-09-29
// @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
// @param         zone    		   string         "e.g. public|dmz.. The empty string is usage default zone, is currently firewalld defualt zone"
// @return        error            error          "Possible errors: ALREADY_ENABLED."
func (c *DbusClientSerivce) PermanentAddPort(port, zone string) (err error) {

	if err = checkPort(port); err != nil {
		return err
	}

	if zone == "" {
		zone = c.GetDefaultZone()
	}

	port, protocol := splitPortProtocol(port)

	if path, err := c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	} else {
		obj := c.Conn.Object(object.INTERFACE, path)
		call := obj.Call(object.CONFIG_ZONE_ADDPORT, dbus.FlagNoAutoStart, port, protocol)
		if call.Err != nil {
			return call.Err
		}
		return nil
	}
}

/*
 * @title         GetPort
 * @description   temporary get a firewalld port list
 * @auth          author           2021-10-05
 * @param         zone             string         "The empty string is usage default zone, is currently firewalld defualt zone."
 *                                                   e.g. public|dmz..
 * @return        []list           Port           "Returns port list of zone."
 * @return        error            error          "Possible errors:
 *                                                      INVALID_ZONE"
 */
func (c *DbusClientSerivce) GetPort(zone string) (list []*Port, err error) {

	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_GETPORTS, dbus.FlagNoAutoStart, zone)

	if call.Err != nil {
		return nil, call.Err
	}
	portList := call.Body[0].([][]string)
	for _, value := range portList {
		list = append(list, &Port{
			Port:     value[0],
			Protocol: value[1],
		})
	}
	return
}

/*
 * @title         PermanentGetPort
 * @description   get Permanent configurtion a firewalld port list.
 * @auth          author           2021-10-05
 * @param         zone             string         "The empty string is usage default zone, is currently firewalld defualt zone"
 *														e.g. public|dmz..
 * @return        []list           Port           "Returns port list of zone."
 * @return        error            error          "Possible errors:"
 * 														INVALID_ZONE
 */
func (c *DbusClientSerivce) PermanentGetPort(zone string) (list []*Port, err error) {

	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return nil, err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_GETPORTS, dbus.FlagNoAutoStart)

	if call.Err != nil {
		return nil, call.Err
	}
	portList := call.Body[0].([][]interface{})

	for _, value := range portList {
		list = append(list, &Port{
			Port:     value[0].(string),
			Protocol: value[1].(string),
		})
	}
	return
}

/*
 * @title         RemovePort
 * @description   temporary delete a firewalld port
 * @auth      	  author           2021-10-05
 * @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
 * @param         zone    		   string         "e.g. public|dmz.. The empty string is usage default zone, is currently firewalld defualt zone"
 * @return        bool             string         "Returns name of zone from which the port was removed."
 * @return        error            error          "Possible errors:
 *                                                      INVALID_ZONE,
 *                                                      INVALID_PORT,
 *                                                      MISSING_PROTOCOL,
 *                                                      INVALID_PROTOCOL,
 *                                                      NOT_ENABLED,
 *                                                      INVALID_COMMAND"
 */
func (c *DbusClientSerivce) RemovePort(port, zone string) (b bool, err error) {

	if err = checkPort(port); err != nil {
		return false, err
	}

	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(port)

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_REMOVEPORT, dbus.FlagNoAutoStart, zone, port, protocol)

	if call.Err != nil {
		return false, call.Err
	}
	return true, nil
}

/*
 * @title         PermanentRemovePort
 * @description   Permanently delete (port, protocol) from list of ports of zone.
 * @auth      	  author           2021-10-05
 * @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
 * @param         zone    		   string         "The empty string is usage default zone, is currently firewalld defualt zone"
 * 														e.g. public|dmz.."
 * @return        bool             string         "Returns name of zone from which the port was removed."
 * @return        error            error          "Possible errors:
 *                                                      NOT_ENABLED"
 */
func (c *DbusClientSerivce) PermanentRemovePort(port, zone string) (b bool, err error) {
	if err = checkPort(port); err != nil {
		return false, err
	}

	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(port)

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return false, err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_REMOVEPORT, dbus.FlagNoAutoStart, port, protocol)

	if call.Err != nil {
		return false, call.Err
	}
	return true, nil
}

/************************************************** Protocol area ***********************************************************/

// @title         AddProtocol
// @description   temporary get a firewalld port list
// @auth      	  author           2021-09-29
// @param         zone    		   string         "e.g. public|dmz.. If zone is empty string, use default zone. "
// @param         protocol         string         "e.g. tcp|udp... The protocol can be any protocol supported by the system."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        zoneName         string         "Returns name of zone to which the protocol was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_PROTOCOL, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) AddProtocol(zone, protocol string, timeout int) (list string, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_ADDPROTOCOL, dbus.FlagNoAutoStart, zone, protocol, timeout)

	if call.Err != nil {
		return "", call.Err
	}
	return call.Body[0].(string), nil
}

/************************************************** service area ***********************************************************/

// @title         AddService
// @description   temporary Add service into zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        zoneName         string         "Returns name of zone to which the service was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) AddService(zone, service string, timeout int) (list string, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_ADDSERVICE, dbus.FlagNoAutoStart, zone, service, timeout)

	if call.Err != nil {
		return "", call.Err
	}
	return call.Body[0].(string), nil
}

// @title         PermanentAddService
// @description   Permanent Add service into zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentAddService(zone, service string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_ADDSERVICE, dbus.FlagNoAutoStart, service)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

// @title         QueryService
// @description   temporary check whether service has been added for zone..
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) QueryService(zone, service string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.PATH)
	call := obj.Call(object.ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, zone, service)
	if !call.Body[0].(bool) {
		return false
	}
	return true
}

// @title         PermanentQueryService
// @description   Permanent Return whether Add service in rich rules in zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentQueryService(zone, service string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	var err error
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return false
	}

	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, service)
	if !call.Body[0].(bool) {
		return false
	}
	return true
}

// @title         RemoveService
// @description   temporary Remove service from zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemoveService(zone, service string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.PATH)
	call := obj.Call(object.ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, zone, service)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

// @title         PermanentAddService
// @description   Permanent Add service into zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentRemoveService(zone, service string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, service)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

/************************************************** Masquerade area ***********************************************************/

/*
 * @title         EnableMasquerade
 * @description   temporary enable masquerade in zone..
 * @auth          author           2021-09-29
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         timeout          int            "Timeout, If timeout is non-zero, masquerading will be active for the amount of seconds."
 * @return        error            error          "Possible errors:
 *                                                  INVALID_ZONE,
 *                                                  ALREADY_ENABLED,
 *                                                  INVALID_COMMAND"
 */
func (c *DbusClientSerivce) EnableMasquerade(zone string, timeout int) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_ADDMASQUERADE, dbus.FlagNoAutoStart, zone, timeout)

	if call.Err != nil && len(call.Body) <= 0 {
		return call.Err
	}
	return
}

/*
 * @title         PermanentEnableMasquerade
 * @description   permanent enable masquerade in zone..
 * @auth          author           2021-09-29
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                  INVALID_ZONE,
 *                                                  ALREADY_ENABLED,
 *                                                  INVALID_COMMAND"
 */
func (c *DbusClientSerivce) PermanentEnableMasquerade(zone string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_ADDMASQUERADE, dbus.FlagNoAutoStart)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

/*
 * @title         DisableMasquerade
 * @description   temporary enable masquerade in zone..
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         timeout          int            "Timeout, If timeout is non-zero, masquerading will be active for the amount of seconds."
 * @return        zoneName         string         "Returns name of zone in which the masquerade was enabled."
 * @return        error            error          "Possible errors:
 *                                                  INVALID_ZONE,
 *                                                  NOT_ENABLED,
 *                                                  INVALID_COMMAND"
 */
func (c *DbusClientSerivce) DisableMasquerade(zone string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_REMOVEMASQUERADE, dbus.FlagNoAutoStart, zone)

	if call.Err != nil && len(call.Body) <= 0 {
		return call.Err
	}
	return
}

/*
 * @title         PermanentDisableMasquerade
 * @description   permanent enable masquerade in zone..
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        b            	   bool           "Possible errors:
 * @return        error            error          "Possible errors:
 *                                                  NOT_ENABLED"
 */
func (c *DbusClientSerivce) PermanentDisableMasquerade(zone string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_REMOVEMASQUERADE, dbus.FlagNoAutoStart)

	if call.Err != nil && len(call.Body) <= 0 {
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentQueryMasquerade
 * @description   query runtime masquerading has been enabled in zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        b            	   bool           "enable: true, disable:false:
 * @return        error            error          "Possible errors:
 *                                                   INVALID_ZONE"
 */
func (c *DbusClientSerivce) PermanentQueryMasquerade(zone string) (b bool, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return false, err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_QUERYMASQUERADE, dbus.FlagNoAutoStart)

	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		return false, call.Err
	}
	return true, nil
}

/*
 * @title         QueryMasquerade
 * @description   query runtime masquerading has been enabled in zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         timeout          int            "Timeout, If timeout is non-zero, masquerading will be active for the amount of seconds."
 * @return        zoneName         string         "Returns name of zone in which the masquerade was enabled."
 * @return        error            error          "Possible errors:
 *                                                  INVALID_ZONE"
 */
func (c *DbusClientSerivce) QueryMasquerade(zone string) (b bool, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_QUERYMASQUERADE, dbus.FlagNoAutoStart, zone)
	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		return false, call.Err
	}
	return true, nil
}

/************************************************** Interface area ***********************************************************/

/*
 * @title         BindInterface
 * @description   temporary Bind interface with zone. From now on all traffic
 * 				   going through the interface will respect the zone's settings.
 * @auth          author           2021-09-29
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        zoneName         string         "Returns name of zone to which the interface was bound."
 * @return        error            error          "Possible errors:
 *                                                      INVALID_ZONE,
 *                                                      INVALID_INTERFACE,
 *                                                      ALREADY_ENABLED,
 *                                                      INVALID_COMMAND"
 */
func (c *DbusClientSerivce) BindInterface(zone, interface_name string) (list string, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)

	if call.Err != nil {
		return "", call.Err
	}
	return call.Body[0].(string), nil
}

/*
 * @title         PermanentBindInterface
 * @description   Permanently Bind interface with zone. From now on all traffic
 * 				   going through the interface will respect the zone's settings.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentBindInterface(zone, interface_name string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, interface_name)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

/*
 * @title         QueryInterface
 * @description   temporary Query whether interface has been bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         interface        string         "device name， e.g. "
 * @return        b         	   bool           "true:enable, fales:disable."
 * @return        error            error          "Possible errors:
 *                                                      INVALID_ZONE,
 *                                                      INVALID_INTERFACE
 */
func (c *DbusClientSerivce) QueryInterface(zone, interface_name string) (b bool, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_QUERYINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)

	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		return false, call.Err
	}
	return true, nil
}

/*
 * @title         PermanentQueryInterface
 * @description   Permanently Query whether interface has been bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentQueryInterface(zone, interface_name string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, interface_name)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

/*
 * @title         RemoveInterface
 * @description   Permanently Query whether interface has been bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) RemoveInterface(zone, interface_name string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_REMOVEINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)
	fmt.Println(call.Body)
	if call.Err != nil && len(call.Body) <= 0 {
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentRemoveInterface
 * @description   Permanently remove interface from list of interfaces bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         interface_name   string         "interface name. e.g. eth0 | ens33.  "
 * @return        error            error          "Possible errors:
 *                                                       NOT_ENABLED"
 */
func (c *DbusClientSerivce) PermanentRemoveInterface(zone, interface_name string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_REMOVEINTERFACE, dbus.FlagNoAutoStart, interface_name)
	fmt.Println(call.Body)
	if call.Err != nil {
		return call.Err
	}
	return nil
}

/************************************************** ForwardPort area ***********************************************************/

/*
 * @title         AddForwardPort
 * @description   temporary Add the IPv4 forward port into zone.
 * @auth      	  author           2021-09-29
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
 * @return        error            error          "Possible errors:
 * 													INVALID_ZONE,
 * 													INVALID_PORT,
 * 													MISSING_PROTOCOL,
 * 													INVALID_PROTOCOL,
 * 													INVALID_ADDR,
 * 													INVALID_FORWARD,
 * 													ALREADY_ENABLED,
 * 													INVALID_COMMAND"
 */
func (c *DbusClientSerivce) AddForwardPort(zone string, portProtocol, toHostPort string, timeout int) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, err := net.SplitHostPort(toHostPort)
	if err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_ADDFORWARDPORT, dbus.FlagNoAutoStart, zone, port, protocol, toPort, toAddr, timeout)
	if call.Err != nil && len(call.Body) <= 0 {
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentAddForwardPort
 * @description   temporary Add the IPv4 forward port into zone.
 * @auth      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentAddForwardPort(zone string, portProtocol, toHostPort string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, err := net.SplitHostPort(toHostPort)
	if err != nil {
		return err
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_ADDFORWARDPORT, dbus.FlagNoAutoStart, port, protocol, toPort, toAddr)
	if call.Err != nil && len(call.Body) <= 0 {
		return call.Err
	}
	return nil
}

/*
 * @title         RemoveForwardPort
 * @description   temporary (runtime) Remove IPv4 forward port ((port, protocol, toport, toaddr)) from zone.
 * @auth      	  author           2021-09-29
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													INVALID_ZONE,
 * 													INVALID_PORT,
 * 													MISSING_PROTOCOL,
 * 													INVALID_PROTOCOL,
 * 													INVALID_ADDR,
 * 													INVALID_FORWARD,
 * 													ALREADY_ENABLED,
 * 													INVALID_COMMAND"
 */
func (c *DbusClientSerivce) RemoveForwardPort(zone string, portProtocol, toHostPort string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, err := net.SplitHostPort(toHostPort)
	if err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_REMOVEFORWARDPORT, dbus.FlagNoAutoStart, zone, port, protocol, toPort, toAddr)
	if call.Err != nil && len(call.Body) <= 0 {
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentRemoveForwardPort
 * @description   Permanently remove (port, protocol, toport, toaddr) from list of forward ports of zone.
 * @auth      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentRemoveForwardPort(zone string, portProtocol, toHostPort string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, err := net.SplitHostPort(toHostPort)
	if err != nil {
		return err
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_REMOVEFORWARDPORT, dbus.FlagNoAutoStart, port, protocol, toPort, toAddr)
	if call.Err != nil && len(call.Body) <= 0 {
		return call.Err
	}
	return nil
}

/*
 * @title         QueryForwardPort
 * @description   temporary (runtime) query whether the IPv4 forward port (port, protocol, toport, toaddr) has been added into zone.
 * @auth      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													INVALID_ZONE,
 * 													INVALID_PORT,
 * 													MISSING_PROTOCOL,
 * 													INVALID_PROTOCOL,
 * 													INVALID_ADDR,
 * 													INVALID_FORWARD,
 * 													ALREADY_ENABLED,
 * 													INVALID_COMMAND"
 */
func (c *DbusClientSerivce) QueryForwardPort(zone string, portProtocol, toHostPort string) (b bool) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, err := net.SplitHostPort(toHostPort)
	if err != nil {
		return false
	}
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_QUERYFORWARDPORT, dbus.FlagNoAutoStart, zone, port, protocol, toPort, toAddr)
	fmt.Println(call.Body)
	if call.Err != nil || !call.Body[0].(bool) {
		return false
	}
	return true
}

/*
 * @title         PermanentQueryForwardPort
 * @description   Permanently remove (port, protocol, toport, toaddr) from list of forward ports of zone.
 * @auth      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentQueryForwardPort(zone string, portProtocol, toHostPort string) (b bool, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, err := net.SplitHostPort(toHostPort)
	if err != nil {
		return false, err
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return false, err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_QUERYFORWARDPORT, dbus.FlagNoAutoStart, port, protocol, toPort, toAddr)
	if call.Err != nil || (len(call.Body) <= 0 || !call.Body[0].(bool)) {
		return false, call.Err
	}
	return true, nil
}

/************************************************** rich rule area ***********************************************************/

// @title         GetRichRules
// @description   Get list of rich-language rules in zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @return        zoneName         string         "Returns name of zone to which the interface was bound."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) GetRichRules(zone string) (ruleList []*Rule, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_GETRICHRULES, dbus.FlagNoAutoStart, zone)

	if call.Err != nil {
		return nil, call.Err
	}
	for _, value := range call.Body[0].([]string) {
		ruleList = append(ruleList, StringToRule(value))
	}

	return
}

// @title         AddRichRule
// @description   temporary Add rich language rule into zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) AddRichRule(zone string, rule *Rule, timeout int) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString(), timeout)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

// @title         PermanentAddRichRule
// @description   Permanently Add rich language rule into zone.
// @auth      	  author           2021-10-05
// @param         zone    	       sting 		  "If zone is empty string, use default zone. e.g. public|dmz..  ""
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) PermanentAddRichRule(zone string, rule *Rule) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, rule.ToString())

	if call.Err != nil {
		return call.Err
	}
	return nil
}

// @title         RemoveRichRule
// @description   temporary Remove rich rule from zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_RULE, NOT_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemoveRichRule(zone string, rule *Rule) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_REOMVERICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())

	if call.Err != nil {
		return call.Err
	}
	return nil
}

// @title         PermanentAddRichRule
// @description   Permanently Add rich language rule into zone.
// @auth      	  author           2021-10-05
// @param         zone    	       sting 		  "If zone is empty string, use default zone. e.g. public|dmz..  ""
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) PermanentRemoveRichRule(zone string, rule *Rule) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_REOMVERICHRULE, dbus.FlagNoAutoStart, rule.ToString())

	if call.Err != nil {
		return call.Err
	}
	return nil
}

// @title         PermanentQueryRichRule
// @description   Check Permanent Configurtion whether rich rule rule has been added in zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        bool             bool           "Possible errors: INVALID_ZONE, INVALID_RULE"
func (c *DbusClientSerivce) PermanentQueryRichRule(zone string, rule *Rule) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	var err error
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return false
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, rule.ToString())

	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		return false
	}
	return true
}

// @title         QueryRichRule
// @description   Check whether rich rule rule has been added in zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        bool             bool           "Possible errors: INVALID_ZONE, INVALID_RULE"
func (c *DbusClientSerivce) QueryRichRule(zone string, rule *Rule) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())

	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		return false
	}
	return true
}

/************************************************** fw service area ***********************************************************/

/*
 * @title         Reload
 * @description   temporary Add rich language rule into zone.
 * @auth          author           2021-10-05
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) Reload() (err error) {
	obj := c.Conn.Object(object.INTERFACE, object.SERVICE)
	call := obj.Call(object.INTERFACE_RELOAD, dbus.FlagNoAutoStart)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

/*
 * @title         flush currently zone zoneSettings to default zoneSettings.
 * @description   temporary Add rich language rule into zone.
 * @auth          author           2021-10-05
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) RuntimeFlush(zone string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	zoneSettings := &Settings{
		Targe:       "default",
		Description: "reset to firewalld-api",
		Short:       "public",
		Interface:   nil,
		Service: []string{
			"ssh",
			"dhcpv6-client",
		},
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, object.ZONE_PATH); err != nil {
		return err
	}
	obj := c.Conn.Object(object.INTERFACE, path)
	call := obj.Call(object.CONFIG_UPDATE, dbus.FlagNoAutoStart, zoneSettings)

	if call.Err != nil || len(call.Body) > 0 {
		return call.Err
	}
	return nil
}
