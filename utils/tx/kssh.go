package tx

// import (
// 	"fmt"
// 	"strconv"

// 	"github.com/srl-labs/containerlab/nodes"
// 	"github.com/srl-labs/containerlab/types"

// 	"github.com/labctl/labctl/utils/tx/kssh"
// )

// type KSSHTx struct { // implement the txs interface
// 	TargetNode *types.NodeConfig
// 	Vars       map[string]interface{}
// 	Config     []configLines
// }

// func (st *KSSHTx) Send(action Action) ([]*Response, error) {

// 	if len(st.Config) == 0 {
// 		return nil, nil
// 	}

// 	var err error
// 	var result []*Response

// 	// Usernames & passwords
// 	ssh_user, ssh_pass := "admin", "admin"
// 	if v, ok := nodes.GetDefaultCredentialsForKind(st.TargetNode.Kind); ok == nil {
// 		ssh_user, ssh_pass = v[0], v[1]
// 	}
// 	ssh_user = get(st.Vars, vkSshUser, ssh_user)
// 	ssh_user = get(st.Vars, vkSshPass, ssh_pass)

// 	// Host & ports, allows an alternative to containerlab's host entries
// 	ssh_host := get(st.Vars, vkSshHost, st.TargetNode.LongName)
// 	ssh_port, err := strconv.Atoi(get(st.Vars, vkSshPort, "22"))
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid port: %s", err)
// 	}

// 	txs, err := kssh.NewSSHTransport(
// 		st.TargetNode,
// 		kssh.WithUserNamePassword(ssh_user, ssh_pass),
// 		kssh.HostKeyCallback(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	txs.Port = ssh_port
// 	txs.Target = st.TargetNode.ShortName

// 	err = txs.Connect(ssh_host)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer txs.Close()

// 	transaction := action != ASend // commit & compare requires transactions

// 	for _, d1 := range st.Config {
// 		//time.Sleep(10 * time.Millisecond)

// 		err := txs.K.Start(txs, transaction)
// 		if err != nil {
// 			return nil, err
// 		}

// 		sshr, err := txs.Write(d1)
// 		if err != nil {
// 			return nil, fmt.Errorf("could not write config: %s\n%s", err, d1)
// 		}

// 		for _, r := range sshr {
// 			if r.Response != "" {
// 				result = append(result, r)
// 			}
// 		}

// 		if transaction {
// 			var resp *Response
// 			//var msg string

// 			switch action {
// 			case ACommit:
// 				resp, err = txs.K.Commit(txs)
// 				// msg = fmt.Sprintf("%s: COMMIT - %d lines", c.TargetNode.ShortName, len(d1))
// 			case ACompare:
// 				resp, err = txs.K.Compare(txs)
// 				if err == nil {
// 					_, err = txs.K.Discard(txs)
// 				}
// 			}
// 			if resp.Response != "" {
// 				result = append(result, resp)
// 				//msg += resp.Slog()
// 			}
// 		}
// 	}
// 	return result, nil
// }
