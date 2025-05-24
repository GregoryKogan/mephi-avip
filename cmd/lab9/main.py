import numpy as np
import soundfile as sf
import matplotlib.pyplot as plt
from scipy import signal

# --- НАСТРОЙКИ ---
# Имя вашего исходного аудиофайла
INPUT_FILENAME = "assets/sample.wav"
# Имена файлов для сохранения результатов
ORIGINAL_SPECTROGRAM_FILENAME = "output/lab9/spectrogram_original.png"
DENOISED_SPECTROGRAM_FILENAME = "output/lab9/spectrogram_denoised.png"
DENOISED_AUDIO_FILENAME = "output/lab9/instrument_denoised.wav"

# Параметры для оконного преобразования Фурье (STFT)
WINDOW_TYPE = "hann"
WINDOW_SIZE_SAMPLES = 2048  # Размер окна в сэмплах
OVERLAP_SAMPLES = 1024  # Размер перекрытия окон в сэмплах

# Параметры для поиска максимальной энергии
TIME_WINDOW_S = 0.1  # Шаг по времени в секундах (Δt)
FREQ_WINDOW_HZ = 50  # Шаг по частоте в герцах (Δf)

# Длительность фрагмента с шумом в начале файла (в секундах)
NOISE_DURATION_S = 1.0


def plot_spectrogram(freqs, times, Sxx, title, filename):
    """Строит и сохраняет спектрограмму."""
    plt.figure(figsize=(12, 6))

    # Используем pcolormesh для отображения. np.abs(Sxx) - магнитуда
    # Добавляем маленькое значение (1e-9), чтобы избежать log(0)
    plt.pcolormesh(times, freqs, 10 * np.log10(np.abs(Sxx) + 1e-9), shading="gouraud")

    plt.ylabel("Частота [Гц]")
    plt.xlabel("Время [с]")
    plt.title(title)

    # Устанавливаем логарифмическую шкалу для оси частот
    plt.yscale("log")
    plt.ylim(20, 20000)  # Ограничиваем диапазон частот для наглядности

    plt.colorbar(label="Интенсивность [дБ]")
    plt.tight_layout()
    plt.savefig(filename)
    print(f"✅ Спектрограмма сохранена в файл: {filename}")
    # plt.show() # Раскомментируйте, если хотите отображать график сразу


def find_max_energy_region(Sxx, freqs, times, fs, time_window_s, freq_window_hz):
    """Находит область с максимальной энергией на спектрограмме."""
    # Энергия пропорциональна квадрату магнитуды
    energy = np.abs(Sxx) ** 2

    # Преобразуем временные и частотные окна из секунд/Гц в индексы массива
    time_step_samples = int(time_window_s * fs / OVERLAP_SAMPLES)
    freq_step_indices = int(freq_window_hz * WINDOW_SIZE_SAMPLES / fs)

    max_energy = 0
    best_time_idx = 0
    best_freq_idx = 0

    # Проходим по спектрограмме с заданным шагом
    for time_idx in range(0, energy.shape[1] - time_step_samples, time_step_samples):
        for freq_idx in range(
            0, energy.shape[0] - freq_step_indices, freq_step_indices
        ):

            # Вычисляем суммарную энергию в текущем окне
            window = energy[
                freq_idx : freq_idx + freq_step_indices,
                time_idx : time_idx + time_step_samples,
            ]
            current_energy = np.sum(window)

            if current_energy > max_energy:
                max_energy = current_energy
                best_time_idx = time_idx
                best_freq_idx = freq_idx

    # Находим центральные время и частоту для найденного окна
    center_time = times[best_time_idx + time_step_samples // 2]
    center_freq = freqs[best_freq_idx + freq_step_indices // 2]

    return center_time, center_freq, max_energy


# --- ОСНОВНОЙ СКРИПТ ---

if __name__ == "__main__":
    try:
        # 1. Загрузка аудиофайла
        audio, fs = sf.read(INPUT_FILENAME)
        print(f"Аудиофайл '{INPUT_FILENAME}' успешно загружен.")
        print(f"Частота дискретизации: {fs} Гц, Количество сэмплов: {len(audio)}")

        # Проверка на моно. Если стерео, берем левый канал.
        if audio.ndim > 1:
            print("Обнаружено стерео, для анализа используется левый канал.")
            audio = audio[:, 0]

    except FileNotFoundError:
        print(
            f"🛑 Ошибка: Файл '{INPUT_FILENAME}' не найден. Убедитесь, что он находится в той же папке, что и скрипт."
        )
        exit()
    except Exception as e:
        print(f"🛑 Произошла ошибка при чтении файла: {e}")
        exit()

    # 2. Построение спектрограммы для оригинального сигнала
    print("\n--- Анализ оригинального сигнала ---")
    freqs, times, Sxx_original = signal.stft(
        audio,
        fs=fs,
        window=WINDOW_TYPE,
        nperseg=WINDOW_SIZE_SAMPLES,
        noverlap=OVERLAP_SAMPLES,
    )
    plot_spectrogram(
        freqs,
        times,
        Sxx_original,
        "Спектрограмма оригинального сигнала",
        ORIGINAL_SPECTROGRAM_FILENAME,
    )

    # 3. Оценка и вычитание шума
    print("\n--- Подавление шума ---")
    # Выделяем фрагмент с шумом из начала записи
    noise_samples = int(NOISE_DURATION_S * fs)
    noise_part = audio[:noise_samples]

    # STFT для шумового фрагмента
    _, _, Sxx_noise = signal.stft(
        noise_part,
        fs=fs,
        window=WINDOW_TYPE,
        nperseg=WINDOW_SIZE_SAMPLES,
        noverlap=OVERLAP_SAMPLES,
    )

    # Оцениваем средний уровень шума по частотам
    mean_noise = np.mean(np.abs(Sxx_noise), axis=1, keepdims=True)

    # Вычитаем спектр шума из спектра сигнала (спектральное вычитание)
    Sxx_denoised_magnitude = np.abs(Sxx_original) - mean_noise
    # Убираем возможные отрицательные значения
    Sxx_denoised_magnitude = np.maximum(Sxx_denoised_magnitude, 0)

    # Восстанавливаем фазу из оригинального сигнала
    phase = np.angle(Sxx_original)
    Sxx_denoised = Sxx_denoised_magnitude * np.exp(1j * phase)

    # Построение спектрограммы для очищенного сигнала
    plot_spectrogram(
        freqs,
        times,
        Sxx_denoised,
        "Спектрограмма после вычитания шума",
        DENOISED_SPECTROGRAM_FILENAME,
    )

    # Восстановление звуковой дорожки и сохранение
    _, audio_denoised = signal.istft(
        Sxx_denoised,
        fs=fs,
        window=WINDOW_TYPE,
        nperseg=WINDOW_SIZE_SAMPLES,
        noverlap=OVERLAP_SAMPLES,
    )
    sf.write(DENOISED_AUDIO_FILENAME, audio_denoised, fs)
    print(f"✅ Очищенный аудиофайл сохранен в: {DENOISED_AUDIO_FILENAME}")
    print("Сравните на слух оригинальный и очищенный файлы.")

    # 4. Поиск момента времени с наибольшей энергией
    print("\n--- Поиск максимальной энергии ---")
    time_res, freq_res, max_e = find_max_energy_region(
        Sxx_original, freqs, times, fs, TIME_WINDOW_S, FREQ_WINDOW_HZ
    )
    print(f"Наибольшая энергия обнаружена в окрестности:")
    print(f"  -> Время: {time_res:.2f} с")
    print(f"  -> Частота: {freq_res:.2f} Гц")
