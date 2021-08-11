module github.com/tapairmax/memcached-operator-with-okt

go 1.16

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0

	// ADDED4OKT
	gitlab.tech.orange/dbmsprivate/operators/okt v1.5.0-beta.3
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	sigs.k8s.io/controller-runtime v0.9.5
)

// ADDED4OKT: To use OKT as a library located on your machine, uncomment the first replace rule and comment the 2nd replace rule.
replace gitlab.tech.orange/dbmsprivate/operators/okt => /home/cloud/go/src/Operators-Karma-Tools

//replace gitlab.tech.orange/dbmsprivate/operators/okt => gitlab.tech.orange/dbmsprivate/operators/okt.git v1.5.0-beta.3
