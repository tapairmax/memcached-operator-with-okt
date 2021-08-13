module github.com/tapairmax/memcached-operator-with-okt

go 1.16

require (

	// ADDED4OKT
	github.com/Orange-OpenSource/Operators-Karma-Tools v1.11.0-beta.2
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	sigs.k8s.io/controller-runtime v0.9.5
)

// ADDED4OKT: To use OKT as a library located on your machine, uncomment the first replace rule and comment the 2nd replace rule.
replace github.com/Orange-OpenSource/Operators-Karma-Tools => /home/cloud/go/src/Operators-Karma-Tools

// No longer needed on Github
//replace gitlab.xxxxx.orange/yyyyy/okt => gitlab.xxxxx.orange/yyyyy/okt.git v1.5.0-beta.3
