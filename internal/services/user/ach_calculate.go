package user

import (
	"context"
	"fmt"

	"gitlab.com/krespix/gamification-api/internal/models"
)

func (s *service) AddEvent(ctx context.Context, email string, eventID int) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	err = s.userRepo.CreateUserEvent(ctx, int(user.ID), eventID)
	if err != nil {
		return err
	}

	achToRecalculate, err := s.getAchievementsToRecalculate(ctx, int(user.ID), 0, eventID)
	if err != nil {
		return err
	}
	fmt.Println(len(achToRecalculate))
	for _, v := range achToRecalculate {
		fmt.Println(v.ID)
	}

	err = s.calculateAchievements(ctx, int(user.ID), achToRecalculate)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddStat(ctx context.Context, email string, statID int) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	err = s.userRepo.CreateUserStat(ctx, int(user.ID), statID)
	if err != nil {
		return err
	}

	achToRecalculate, err := s.getAchievementsToRecalculate(ctx, int(user.ID), statID, 0)
	if err != nil {
		return err
	}
	fmt.Println(len(achToRecalculate))
	for _, v := range achToRecalculate {
		fmt.Println(v.ID)
	}

	err = s.calculateAchievements(ctx, int(user.ID), achToRecalculate)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) calculateAchievements(ctx context.Context, userID int, achsList []*models.RepoAchievement) error {
	for _, a := range achsList {
		achieved, err := s.calculateSingleAchievement(ctx, userID, a)
		if err != nil {
			return err
		}
		if achieved {
			err = s.achievementsRepo.CreateUserAch(ctx, userID, a.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) calculateSingleAchievement(ctx context.Context, userID int, ach *models.RepoAchievement) (bool, error) {
	userEvents, err := s.eventsRepo.GetUserEvents(ctx, userID)
	if err != nil {
		return false, err
	}
	userEventsMap := make(map[int]*models.UserEvent, len(userEvents))
	for i, v := range userEvents {
		userEventsMap[v.EventID] = userEvents[i]
	}
	userStats, err := s.statRepo.GetUserStats(ctx, userID)
	if err != nil {
		return false, err
	}
	userStatsMap := make(map[int]*models.UserStat, len(userStats))
	for i, v := range userStats {
		userStatsMap[v.StatID] = userStats[i]
	}

	blocksRes := make([]int, 0, len(ach.Rules.Blocks))

	for _, block := range ach.Rules.Blocks {
		res, err := s.calculateBlock(ctx, block, userEventsMap, userStatsMap)
		if err != nil {
			return false, err
		}
		blocksRes = append(blocksRes, res)
	}

	sum := 0
	for _, v := range blocksRes {
		sum += v
	}

	if sum > 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) calculateBlock(ctx context.Context, block *models.RulesBlock, eventsMap map[int]*models.UserEvent, statsMap map[int]*models.UserStat) (int, error) {
	resBlock := make([]int, 0)
	for _, eventRule := range block.EventRules {
		if eventRule.NeedParticipate {
			if _, ok := eventsMap[eventRule.EventID]; ok {
				resBlock = append(resBlock, 1)
			} else {
				return 0, nil
			}
		} else {
			if _, ok := eventsMap[eventRule.EventID]; ok {
				return 0, nil
			} else {
				resBlock = append(resBlock, 1)
			}
		}
	}

	for _, statRule := range block.StatRules {
		if _, ok := statsMap[statRule.StatID]; !ok {
			return 0, nil
		}
		isDone := s.compareStat(statsMap[statRule.StatID].Value, statRule.TargetValue, statRule.Comparison)
		if isDone {
			resBlock = append(resBlock, 1)
		} else {
			return 0, nil
		}
	}

	res := 1
	for _, v := range resBlock {
		res *= v
	}

	return res, nil
}

func (s *service) compareStat(userValue, targetValue int, comparison models.Comparison) bool {
	switch comparison {
	case models.ComparisonEquals:
		return userValue == targetValue
	case models.ComparisonLesserThan:
		return userValue < targetValue
	case models.ComparisonGreaterThan:
		return userValue > targetValue
	case models.ComparisonNotEquals:
		return userValue == targetValue
	default:
		return false
	}
}

func (s *service) getAchievementsToRecalculate(ctx context.Context, userID int, statID int, eventID int) ([]*models.RepoAchievement, error) {
	userAchs, err := s.achievementsRepo.GetUsersAchievements(ctx, userID)
	if err != nil {
		return nil, err
	}
	userAchsMap := make(map[int]*models.UserAch, len(userAchs))
	for _, ua := range userAchs {
		userAchsMap[ua.AchID] = &models.UserAch{
			AchID:       ua.AchID,
			CreatedAt:   ua.CreatedAt,
			Name:        ua.Name,
			Description: ua.Description,
			Image:       ua.Image,
		}
	}

	allAchs, err := s.achievementsRepo.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	achToRecalculate := make([]*models.RepoAchievement, 0)
	recalculateMap := make(map[int]*models.RepoAchievement, 0)

	for _, ach := range allAchs {
		_, ok := userAchsMap[ach.ID]
		_, alreadyExists := recalculateMap[ach.ID]
		if alreadyExists {
			continue
		}
		for _, b := range ach.Rules.Blocks {
			for _, sr := range b.StatRules {
				if sr.StatID == statID && !ok {
					recalculateMap[ach.ID] = &models.RepoAchievement{
						ID:          ach.ID,
						Name:        ach.Name,
						Description: ach.Description,
						Image:       ach.Image,
						EndAt:       ach.EndAt,
						CreatedAt:   ach.CreatedAt,
						Rules:       ach.Rules,
					}
				}
			}
			for _, er := range b.EventRules {
				if er.EventID == eventID && !ok {
					recalculateMap[ach.ID] = &models.RepoAchievement{
						ID:          ach.ID,
						Name:        ach.Name,
						Description: ach.Description,
						Image:       ach.Image,
						EndAt:       ach.EndAt,
						CreatedAt:   ach.CreatedAt,
						Rules:       ach.Rules,
					}
				}
			}
		}

	}

	for k := range recalculateMap {
		achToRecalculate = append(achToRecalculate, recalculateMap[k])
	}

	return achToRecalculate, nil
}
