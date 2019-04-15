package cn_proto

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tianhphahai2/gRPC-CRUD/pkg/api/cn_proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "cn_proto"
)

// test_gRPCServiceServer
type test_gRPCServiceServer struct {
	db *sql.DB
}

// connect SQL database
func (s *test_gRPCServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create
func (s *test_gRPCServiceServer) Create(ctx context.Context, req *cn_proto.CreateRequest) (*cn_proto.CreateResponse, error) {
	// check API
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// connect db
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// insert from db
	res, err := c.ExecContext(ctx, "insert into test__grpc(`Title`, `Descript`) values (?,?)",
		req.TestGrpc.Title, req.TestGrpc.Description)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into test_gRPC-> "+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created test_gRPC-> "+err.Error())
	}

	return &cn_proto.CreateResponse{
		Api:apiVersion,
		Id:id,
	}, nil
}

// read
func (s *test_gRPCServiceServer) Read(ctx context.Context, req *cn_proto.ReadRequest) (*cn_proto.ReadResponse, error) {
	// check api
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get sql
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query test_grpc by ID
	rows, err := c.QueryContext(ctx, "SELECT `id`, `Title`, `Descript` from test__grpc where `id`=?",
		req.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from test_gRPC-> "+err.Error())
	}

	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from test_gRPC-> "+err.Error())
		}
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("test_grpc with ID='%d' is not found",
			req.Id))
	}

	// get test_rgpc data
	var td cn_proto.TestGrpc

	if err := rows.Scan(&td.Id, &td.Title, &td.Description); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from test_grpc row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("found multiple test_grpc rows with ID='%d'",
			req.Id))
	}

	return &cn_proto.ReadResponse{
		Api:apiVersion,
		TestGrpc: &td,
	}, nil
}

// Delete
func (s *test_gRPCServiceServer) Delete(ctx context.Context, req *cn_proto.DeleteRequest) (*cn_proto.DeleteResponse, error) {
	// check api
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// connect DB
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// delete test_grpc
	res, err := c.ExecContext(ctx, "DELETE FROM test__grpc where `id`=?",
		req.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete test_grpc-> "+err.Error())
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.Id))
	}

	return &cn_proto.DeleteResponse{
		Api:apiVersion,
		Deleted:rows,
	}, nil
}

// update
func (s *test_gRPCServiceServer) Update(ctx context.Context, req *cn_proto.UpdateRequest) (*cn_proto.UpdateResponse, error) {
	// check api
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// connect db
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// update test_grpc
	res, err := c.ExecContext(ctx, "UPDATE test__grpc SET `Title`=?. `Descript`=? WHERE `id`=?",
		req.TestGrpc.Title, req.TestGrpc.Description)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update test_grpc-> "+err.Error())
	}

	rows, err := res.RowsAffected()

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("test_grpc with ID='%d' is not found", req.TestGrpc.Id))
	}

	return &cn_proto.UpdateResponse{
		Api:apiVersion,
		Updated:rows,
	}, nil
}

// read All
func (s *test_gRPCServiceServer) ReadAll(ctx context.Context, req *cn_proto.ReadAllRequest) (*cn_proto.ReadAllResponse, error) {
	// check api
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// connect db
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get db
	rows, err := c.QueryContext(ctx, "select `id`, `Title`, `Descript` from test__grpc")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from test_grpc-> "+err.Error())
	}
	defer rows.Close()

	list := []*cn_proto.TestGrpc{}
	if rows.Next() {
		td := new(cn_proto.TestGrpc)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from test_grpc row-> "+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from test_grpc-> "+err.Error())
	}

	return &cn_proto.ReadAllResponse{
		Api:apiVersion,
		TestGrpc: list,
	}, nil
}

// New db create test_gRPC service
func NewTest_gRPCServiceServer(db *sql.DB) cn_proto.TestGrpcServiceServer  {
	return &test_gRPCServiceServer{db: db}
}

// func check API
func (s *test_gRPCServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}