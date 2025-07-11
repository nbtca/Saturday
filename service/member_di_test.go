package service

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
)

// Mock implementations for testing dependency injection

type mockMemberRepository struct {
	members map[string]model.Member
}

func newMockMemberRepository() *mockMemberRepository {
	return &mockMemberRepository{
		members: make(map[string]model.Member),
	}
}

func (m *mockMemberRepository) GetMemberById(id string) (model.Member, error) {
	if member, exists := m.members[id]; exists {
		return member, nil
	}
	return model.Member{}, nil
}

func (m *mockMemberRepository) GetMemberByLogtoId(logtoId string) (model.Member, error) {
	for _, member := range m.members {
		if member.LogtoId == logtoId {
			return member, nil
		}
	}
	return model.Member{}, nil
}

func (m *mockMemberRepository) GetMemberByGithubId(githubId string) (model.Member, error) {
	for _, member := range m.members {
		if member.GithubId == githubId {
			return member, nil
		}
	}
	return model.Member{}, nil
}

func (m *mockMemberRepository) GetMemberIdByLogtoId(logtoId string) (sql.NullString, error) {
	for id, member := range m.members {
		if member.LogtoId == logtoId {
			return sql.NullString{String: id, Valid: true}, nil
		}
	}
	return sql.NullString{Valid: false}, nil
}

func (m *mockMemberRepository) GetMembers(offset uint64, limit uint64) ([]model.Member, error) {
	var result []model.Member
	for _, member := range m.members {
		result = append(result, member)
	}
	return result, nil
}

func (m *mockMemberRepository) CreateMember(member *model.Member) error {
	m.members[member.MemberId] = *member
	return nil
}

func (m *mockMemberRepository) UpdateMember(member model.Member) error {
	m.members[member.MemberId] = member
	return nil
}

func (m *mockMemberRepository) ExistMember(id string) (bool, error) {
	_, exists := m.members[id]
	return exists, nil
}

type mockRoleRepository struct{}

func newMockRoleRepository() *mockRoleRepository {
	return &mockRoleRepository{}
}

func (m *mockRoleRepository) ExistRole(role string) (bool, error) {
	validRoles := []string{"admin", "member", "admin_inactive", "member_inactive"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockRoleRepository) SetMemberRole(memberId string, role string, conn *sqlx.Tx) error {
	return nil
}

type mockLogtoService struct{}

func newMockLogtoService() *mockLogtoService {
	return &mockLogtoService{}
}

func (m *mockLogtoService) FetchLogtoToken(resource string, scope string) (map[string]interface{}, error) {
	return map[string]interface{}{"access_token": "mock_token"}, nil
}

func (m *mockLogtoService) FetchUsers(request FetchLogtoUsersRequest) ([]FetchLogtoUsersResponse, error) {
	return []FetchLogtoUsersResponse{}, nil
}

func (m *mockLogtoService) FetchUserById(userId string) (*FetchLogtoUsersResponse, error) {
	return &FetchLogtoUsersResponse{Id: userId, Name: "Test User"}, nil
}

func (m *mockLogtoService) FetchUserByToken(token string) (*FetchLogtoUsersResponse, error) {
	return &FetchLogtoUsersResponse{Id: "test-user-id", Name: "Test User"}, nil
}

func (m *mockLogtoService) FetchUserInfo(accessToken string) (FetchUserInfoResponse, error) {
	return FetchUserInfoResponse{}, nil
}

func (m *mockLogtoService) PatchUserById(userId string, data dto.PatchLogtoUserRequest) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (m *mockLogtoService) FetchUserRole(userId string) (FetchUserRoleResponse, error) {
	return []LogtoUserRole{
		{Name: "Repair Member", Type: "test"},
	}, nil
}

func TestMemberServiceDependencyInjection(t *testing.T) {
	// Create mock dependencies
	memberRepo := newMockMemberRepository()
	roleRepo := newMockRoleRepository()
	logtoService := newMockLogtoService()
	
	// Create service with injected dependencies
	memberService := NewMemberService(memberRepo, roleRepo, logtoService)
	
	// Test service functionality
	testMember := &model.Member{
		MemberId: "test123",
		Alias:    "Test User",
		Role:     "member",
		LogtoId:  "logto-123",
	}
	
	// Test CreateMember
	err := memberService.CreateMember(testMember)
	if err != nil {
		t.Errorf("CreateMember failed: %v", err)
	}
	
	// Test GetMemberById
	retrievedMember, err := memberService.GetMemberById("test123")
	if err != nil {
		t.Errorf("GetMemberById failed: %v", err)
	}
	
	if retrievedMember.MemberId != testMember.MemberId {
		t.Errorf("Expected MemberId %s, got %s", testMember.MemberId, retrievedMember.MemberId)
	}
	
	// Test GetMemberByLogtoId
	retrievedByLogto, err := memberService.GetMemberByLogtoId("logto-123")
	if err != nil {
		t.Errorf("GetMemberByLogtoId failed: %v", err)
	}
	
	if retrievedByLogto.LogtoId != testMember.LogtoId {
		t.Errorf("Expected LogtoId %s, got %s", testMember.LogtoId, retrievedByLogto.LogtoId)
	}
}

func TestMemberServiceMapLogtoUserRole(t *testing.T) {
	// Create mock dependencies
	memberRepo := newMockMemberRepository()
	roleRepo := newMockRoleRepository()
	logtoService := newMockLogtoService()
	
	// Create service with injected dependencies
	memberService := NewMemberService(memberRepo, roleRepo, logtoService)
	
	// Cast to concrete type to access MapLogtoUserRole method
	concreteMemberService, ok := memberService.(*MemberService)
	if !ok {
		t.Fatal("Failed to cast to concrete MemberService type")
	}
	
	// Test role mapping
	testCases := []struct {
		roles    []LogtoUserRole
		expected string
	}{
		{
			roles:    []LogtoUserRole{{Name: "Repair Admin"}},
			expected: "admin",
		},
		{
			roles:    []LogtoUserRole{{Name: "Repair Member"}},
			expected: "member",
		},
		{
			roles:    []LogtoUserRole{{Name: "Repair Admin"}, {Name: "Repair Member"}},
			expected: "admin", // Admin takes precedence
		},
		{
			roles:    []LogtoUserRole{{Name: "Unknown Role"}},
			expected: "",
		},
	}
	
	for _, tc := range testCases {
		result := concreteMemberService.MapLogtoUserRole(tc.roles)
		if result != tc.expected {
			t.Errorf("Expected role %s, got %s for roles %v", tc.expected, result, tc.roles)
		}
	}
}