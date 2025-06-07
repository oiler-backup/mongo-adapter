package restorer

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	ctx        = context.Background()
	dbUser     = "root"
	dbPass     = "pass"
	dbName     = "admin"
	backupName = "backup.dump"
)

func setupMongoContainer() (*tc.Container, error) {
	req := tc.ContainerRequest{
		Image:           "mongo:8.0",
		ExposedPorts:    []string{"27017/tcp"},
		AlwaysPullImage: false,
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": dbUser,
			"MONGO_INITDB_ROOT_PASSWORD": dbPass,
		},
		WaitingFor: wait.ForListeningPort("27017/tcp"),
	}

	mongoC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	mongoC.Exec(ctx, []string{"apt-get mongodb-tools"})
	return &mongoC, err
}

// func Test_Redtore_UploadValidDump(t *testing.T) {
// 	mongoC, err := setupMongoContainer()
// 	require.NoError(t, err)
// 	defer func() {
// 		err := (*mongoC).Terminate(ctx)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	dbhost, _ := (*mongoC).ContainerIP(ctx)
// 	dbPort, _ := (*mongoC).MappedPort(ctx, "27017")
// 	tempDir := t.TempDir()
// 	backupFile := filepath.Join(tempDir, "backup.dump")
// 	file, err := os.Create(backupFile)
// 	if err != nil {
// 		panic(err)
// 	}
// 	file.Close()

// 	r := NewRestorer(
// 		dbhost,
// 		dbPort.Port(),
// 		dbUser,
// 		dbPass,
// 		dbName,
// 		backupFile,
// 	)

// 	err = r.Restore(ctx)
// 	require.NoError(t, err)
// }

func Test_Redtore_InvalidDump(t *testing.T) {
	mongoC, err := setupMongoContainer()
	require.NoError(t, err)
	defer func() {
		err := (*mongoC).Terminate(ctx)
		if err != nil {
			panic(err)
		}
	}()

	dbhost, _ := (*mongoC).ContainerIP(ctx)
	tempDir := t.TempDir()
	backupFile := filepath.Join(tempDir, backupName)

	r := NewRestorer(
		dbhost,
		"3306",
		dbUser,
		dbPass,
		dbName,
		backupFile,
	)

	err = r.Restore(ctx)
	require.ErrorContains(t, err, "failed executing mongorestore:")
}

func Test_Redtore_InvalidDBHost(t *testing.T) {
	dbhost := "wrong"
	dbPort := "3306"
	r := NewRestorer(
		dbhost,
		dbPort,
		dbUser,
		dbPass,
		dbName,
		backupName,
	)

	err := r.Restore(ctx)
	require.ErrorContains(t, err, "failed executing mongorestore:")
}
