ALTER TABLE meal_plans
    ADD CONSTRAINT unique_user_day_weekstart_mealtype
        UNIQUE (user_id, day_of_week, week_start_date, meal_type);
