package dns

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/arshzip/txtllm/internal/llm"
	"github.com/miekg/dns"
)

func Start(client *llm.OpenRouterClient) {
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		for _, q := range r.Question {
			if q.Qtype == dns.TypeTXT || q.Qtype == dns.TypeA {
				prompt := strings.TrimSuffix(q.Name, ".")
				answer, err := client.Query(llm.DefaultModel, prompt)
				if err != nil {
					log.Printf("LLM query failed: %v", err)
					answer = "Error: " + err.Error()
				}

				// idk why web openrouter models love to list references despite being instructed not to
				// remove markdown links of the form []()
				answer = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`).ReplaceAllString(answer, "")

				// remove raw URLs
				answer = regexp.MustCompile(`https?://[^\s]+`).ReplaceAllString(answer, "")

				cleanAnswer := strings.ReplaceAll(answer, "\"", "'")
				rrString := fmt.Sprintf("3600 IN TXT \"%s\"", cleanAnswer)

				rr, err := dns.NewRR(rrString)
				if err != nil {
					log.Printf("Failed to create RR: %v", err)
					continue
				}
				m.Answer = append(m.Answer, rr)
			}
		}
		w.WriteMsg(m)
	})

	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
