package resolvers

import (
	"context"
	"time"

	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
)

func (r *Resolver) CreateAchievement(ctx context.Context, achievement apiModels.CreateAchievement) (interface{}, error) {
	createAch := &models.CreateAchievement{
		Name:  achievement.Name,
		Rules: convertInputRulesToService(achievement.Rules),
		Image: achievement.Image,
	}
	if achievement.EndAt != nil {
		t := *achievement.EndAt
		createAch.EndAt = time.Unix(int64(t), 0)
	}
	if achievement.Description != nil {
		createAch.Description = *achievement.Description
	}

	err := r.achievementService.Create(ctx, createAch)
	return nil, err
}

func (r *Resolver) UpdateAchievement(ctx context.Context, achievement apiModels.UpdateAchievement) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) DeleteAchievement(ctx context.Context, id int) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) GetAchievements(ctx context.Context, pagination *apiModels.Pagination) (*apiModels.GetAchievementsResponse, error) {

	return nil, nil
}

func (r *Resolver) GetAchievement(ctx context.Context, id int) (*apiModels.Achievement, error) {
	ach, err := r.achievementService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return repoAchToAPI(ach), nil
}

func repoAchToAPI(achievement *models.RepoAchievement) *apiModels.Achievement {
	res := &apiModels.Achievement{
		ID:          achievement.ID,
		Name:        achievement.Name,
		Description: utils.SqlNullStringToString(achievement.Description),
		Image:       utils.SqlNullStringToString(achievement.Image),
		CreatedAt:   int(achievement.CreatedAt.Unix()),
	}
	if achievement.EndAt.Valid {
		t := utils.SqlNullTimeToTime(achievement.EndAt)
		res.EndAt = &t
	}
	if achievement.Rules != nil {
		res.Rules = convertServiceRulesToAPI(achievement.Rules)
	}
	return res
}

func convertServiceRulesToAPI(rules *models.Rules) *apiModels.Rules {
	res := &apiModels.Rules{
		Blocks: make([]*apiModels.RuleBlock, 0, len(rules.Blocks)),
	}

	for _, block := range rules.Blocks {
		resBlock := &apiModels.RuleBlock{
			EventsRules: make([]*apiModels.EventRule, 0, len(block.EventRules)),
			StatRules:   make([]*apiModels.StatRule, 0, len(block.StatRules)),
		}
		operator := mapConnectionOperatorToAPI(block.ConnectionOperator)
		resBlock.ConnectionOperator = &operator

		for _, eventRule := range block.EventRules {
			resBlock.EventsRules = append(resBlock.EventsRules, &apiModels.EventRule{
				EventID:         eventRule.EventID,
				NeedParticipate: eventRule.NeedParticipate,
			})
		}

		for _, statRule := range block.StatRules {
			resBlock.StatRules = append(resBlock.StatRules, &apiModels.StatRule{
				StatID:         statRule.StatID,
				TargetValue:    statRule.TargetValue,
				ComparisonType: mapComparisonTypeToAPI(statRule.Comparison),
			})
		}

		res.Blocks = append(res.Blocks, resBlock)
	}

	return res
}

func mapComparisonTypeToAPI(comparison models.Comparison) apiModels.Comparison {
	switch comparison {
	case models.ComparisonNotEquals:
		return apiModels.ComparisonNotEquals
	case models.ComparisonGreaterThan:
		return apiModels.ComparisonGreaterThan
	case models.ComparisonLesserThan:
		return apiModels.ComparisonLesserThan
	case models.ComparisonEquals:
		return apiModels.ComparisonEquals
	default:
		return apiModels.ComparisonInvalidComparison
	}
}

func mapConnectionOperatorToAPI(operator models.ConnectionOperator) apiModels.ConnectionOperator {
	switch operator {
	case models.ConnectionOperatorOr:
		return apiModels.ConnectionOperatorOr
	case models.ConnectionOperatorAnd:
		return apiModels.ConnectionOperatorAnd
	default:
		return apiModels.ConnectionOperatorInvalidConnectionOperator
	}
}

func convertInputRulesToService(rules *apiModels.InputRules) *models.Rules {
	res := &models.Rules{
		Blocks: make([]*models.RulesBlock, 0, len(rules.Blocks)),
	}

	for _, block := range rules.Blocks {
		resBlock := &models.RulesBlock{
			EventRules:         make([]*models.EventRule, 0, len(block.EventsRules)),
			StatRules:          make([]*models.StatRule, 0, len(block.StatRules)),
			ConnectionOperator: mapConnectionOperatorToService(block.ConnectionOperator),
		}
		for _, statRule := range block.StatRules {
			resBlock.StatRules = append(resBlock.StatRules, &models.StatRule{
				StatID:      statRule.StatID,
				TargetValue: statRule.TargetValue,
				Comparison:  mapComparisonTypeToService(statRule.ComparisonType),
			})
		}
		for _, eventRule := range block.EventsRules {
			resBlock.EventRules = append(resBlock.EventRules, &models.EventRule{
				EventID:         eventRule.EventID,
				NeedParticipate: eventRule.NeedParticipate,
			})
		}
		res.Blocks = append(res.Blocks, resBlock)
	}

	return res
}

func mapComparisonTypeToService(comparison apiModels.Comparison) models.Comparison {
	switch comparison {
	case apiModels.ComparisonEquals:
		return models.ComparisonEquals
	case apiModels.ComparisonNotEquals:
		return models.ComparisonNotEquals
	case apiModels.ComparisonGreaterThan:
		return models.ComparisonGreaterThan
	case apiModels.ComparisonLesserThan:
		return models.ComparisonLesserThan
	default:
		return models.ComparisonInvalid
	}
}

func mapConnectionOperatorToService(operator *apiModels.ConnectionOperator) models.ConnectionOperator {
	if operator == nil {
		return models.ConnectionOperatorInvalid
	}
	switch *operator {
	case apiModels.ConnectionOperatorOr:
		return models.ConnectionOperatorOr
	case apiModels.ConnectionOperatorAnd:
		return models.ConnectionOperatorAnd
	default:
		return models.ConnectionOperatorInvalid
	}
}
